package sampler

import (
	"github.com/nareix/joy4/av/avutil"
	"fmt"
	"github.com/dgmann/joy4/cgo/ffmpeg"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/format"
	"github.com/nareix/joy4/av/pktque"
	"time"
	"github.com/dgmann/ma-shared"
)

func init() {
	format.RegisterAll()
}

func Sample(path string) (chan shared.VideoSample) {
	samples := extractSample(path)

	return samples
}

func extractSample(path string) (chan shared.VideoSample) {
	samples := make(chan shared.VideoSample, 10000)

	go func() {
		file, err := avutil.Open(path)
		demuxer := &pktque.FilterDemuxer{Demuxer: file, Filter: &pktque.Walltime{}}
		if err != nil {
			fmt.Errorf("Error %s", err)
		}

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
						Raw:frame.Image,
						FrameNumber: frameCount,
						ReadPacketAt:  readAt,
						CreatedAt: time.Now(),
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
