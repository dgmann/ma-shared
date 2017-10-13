package sampler

import (
	"github.com/dgmann/joy4/cgo/ffmpeg"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/format"
	"github.com/nareix/joy4/av/pktque"
	"time"
	"github.com/dgmann/ma-shared"
	"sync"
)

func init() {
	format.RegisterAll()
}

type VideoSource interface {
	Open() av.Demuxer
}

type VideoServer interface {
	Demuxers() <-chan av.Demuxer
	Listen()
}

type SampleFactory struct {
	source VideoSource
	numSamplers int
}

func NewSampleFactory(source VideoSource, numSamplers int) SampleFactory {
	return SampleFactory{source, numSamplers}
}

func(factory *SampleFactory) StartSampler() chan shared.Message {
	file := factory.source.Open()
	samples := extractSample(file)
	output := make(chan shared.Message, 10000)

	var wg sync.WaitGroup

	for i:=0; i < factory.numSamplers; i++ {
		wg.Add(1)
		go func() {
			for sample := range samples {
				msg, _ := shared.NewMessageFromSample(sample)
				output <- *msg
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(output)
	}()
	return output
}

func Sample(source VideoSource) (chan shared.VideoSample) {
	file := source.Open()
	samples := extractSample(file)

	return samples
}

func SampleMany(source VideoServer) (chan shared.VideoSample) {
	result := make(chan shared.VideoSample, 10000)

	go func() {
		for demuxer := range source.Demuxers() {
			samples := extractSample(demuxer)
			go func() {
				for sample := range samples {
					result <- sample
				}
			}()
		}
	}()
	return result
}

func extractSample(demuxer av.Demuxer) (chan shared.VideoSample) {
	samples := make(chan shared.VideoSample, 10000)

	go func() {

		demuxer := &pktque.FilterDemuxer{Demuxer: demuxer, Filter: &pktque.Walltime{}}

		streams, _ := demuxer.Streams()
		var dec *ffmpeg.VideoDecoder

		for _, stream := range streams {
			if stream.Type().IsVideo() {
				dec, _ = ffmpeg.NewVideoDecoder(stream.(av.VideoCodecData))
			}
		}

		frameCount := 0
		for {
			readAt := time.Now()
			pkt, err := demuxer.ReadPacket()
			if err != nil {
				close(samples)
				return
			}
			if streams[pkt.Idx].Type().IsVideo() {
				frame, _ := dec.Decode(pkt.Data)
				if frame != nil {
					sample := shared.VideoSample{
						Raw:          frame.Image,
						FrameNumber:  frameCount,
						ReadPacketAt: readAt,
						CreatedAt:    time.Now(),
					}
					frame.Free()
					samples <- sample
					frameCount++
				}
			}
		}
	}()

	return samples
}
