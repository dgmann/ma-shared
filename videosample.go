package shared

import (
	"image"
	"time"
	"bytes"
	"bufio"
	"image/jpeg"
)

type VideoSample struct {
	Raw image.YCbCr
	FrameNumber int
	ReadPacketAt time.Time
	CreatedAt time.Time
}

func(sample *VideoSample) ToJPEG() (bytes.Buffer, error) {
	img := sample.Raw
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	err := jpeg.Encode(writer, &img, &jpeg.Options{Quality: 10})
	if err != nil {
		return b, err
	}
	writer.Flush()
	return b, nil
}
