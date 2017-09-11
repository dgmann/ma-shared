package sampler

import (
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av"
)

type FileSource struct {
	Path string
}

func CreateFileSource(path string) *FileSource {
	return &FileSource{Path:path}
}

func(source *FileSource) Open() av.Demuxer {
	file, _ := avutil.Open(source.Path)
	return file
}

func(source *FileSource) Start() {

}

