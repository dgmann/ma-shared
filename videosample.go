package shared

import (
	"image"
	"time"
	"bytes"
	"bufio"
	"image/jpeg"
	"golang.org/x/image/bmp"
)

type VideoSample struct {
	Raw image.YCbCr
	FrameNumber int
	ReadPacketAt time.Time
	CreatedAt time.Time
}

func(sample *VideoSample) ToJPEG() (bytes.Buffer) {
	img := sample.Raw
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	jpeg.Encode(writer, &img, &jpeg.Options{Quality: 10})
	writer.Flush()
	return b
}

func(sample *VideoSample) ToBitmap() (bytes.Buffer) {
	img := sample.Raw
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	bmp.Encode(writer, &img)
	writer.Flush()
	return b
}
