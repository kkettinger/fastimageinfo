package parser

import (
	"bytes"
	"encoding/binary"
)

type BMPParser struct{}

func (B BMPParser) Type() ImageType {
	return BMP
}

func (B BMPParser) DetectType(p []byte) (r Result) {
	bmpHeader := []byte{'B', 'M'}

	if len(p) < len(bmpHeader) {
		return NeedMoreData
	}

	if bytes.Equal(p[0:len(bmpHeader)], bmpHeader) {
		return Valid
	} else {
		return Invalid
	}
}

func (B BMPParser) GetSize(p []byte) (r Result, t ImageSize) {
	if result := B.DetectType(p); result != Valid {
		return result, ImageSize{}
	}

	if len(p) < 26 {
		return NeedMoreData, ImageSize{}
	}

	imageSize := ImageSize{}

	imageSize.Width = binary.LittleEndian.Uint32(p[18:])
	imageSize.Height = binary.LittleEndian.Uint32(p[22:])

	return Valid, imageSize
}

func init() {
	register(&BMPParser{})
}
