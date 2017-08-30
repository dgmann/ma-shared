package sampler

import (
	"github.com/nareix/joy4/av/avutil"
	"fmt"
	"github.com/dgmann/joy4/cgo/ffmpeg"
	"github.com/nareix/joy4/av"
	"image/jpeg"
	"bufio"
	"github.com/nareix/joy4/format"
	"image"
	"github.com/nareix/joy4/av/pktque"
	"time"
	"bytes"
)

func init() {
	format.RegisterAll()
}

type VideoSample struct {
	Raw image.YCbCr
	FrameNumber int
	CreatedAt time.Time
}

func Sample(path string) (chan VideoSample) {
	samples := extractSample(path)

	return samples
}

func extractSample(path string) (chan VideoSample) {
	samples := make(chan VideoSample)

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
			pkt, err := demuxer.ReadPacket()
			if err != nil {
				close(samples)
				return
			}
			if streams[pkt.Idx].Type().IsVideo() {
				frame, _ := dec.Decode(pkt.Data)
				if frame != nil {
					sample := VideoSample{
						Raw:frame.Image,
						FrameNumber: frameCount,
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

func(sample *VideoSample) ToJPEG() (bytes.Buffer) {
	img := sample.Raw
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	jpeg.Encode(writer, &img, &jpeg.Options{Quality: 10})
	writer.Flush()
	return b
}