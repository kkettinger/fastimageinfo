package parser

import (
	"bytes"
	"encoding/binary"
)

type GIFParser struct{}

func (G GIFParser) Type() ImageType {
	return GIF
}

func (G GIFParser) DetectType(p []byte) (r Result) {
	gifHeader := []byte{'G', 'I', 'F'}

	if len(p) < len(gifHeader) {
		return NeedMoreData
	}

	if bytes.Equal(p[0:len(gifHeader)], gifHeader) {
		return Valid
	} else {
		return Invalid
	}
}

func (G GIFParser) GetSize(p []byte) (r Result, t ImageSize) {
	if result := G.DetectType(p); result != Valid {
		return result, ImageSize{}
	}

	if len(p) < 12 {
		return NeedMoreData, ImageSize{}
	}

	imageSize := ImageSize{}
	imageSize.Width = uint32(binary.LittleEndian.Uint16(p[6:]))
	imageSize.Height = uint32(binary.LittleEndian.Uint16(p[8:]))

	return Valid, imageSize
}

func init() {
	register(&GIFParser{})
}
