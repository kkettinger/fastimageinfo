package parser

import (
	"bytes"
	"encoding/binary"
)

type PNGParser struct{}

func (P PNGParser) Type() ImageType {
	return PNG
}

func (P PNGParser) DetectType(p []byte) (r Result) {

	pngFileSignature := []byte{'\x89', 'P', 'N', 'G', '\x0D', '\x0A', '\x1A', '\x0A'}

	if len(p) < len(pngFileSignature) {
		return NeedMoreData
	}

	if bytes.Equal(p[0:len(pngFileSignature)], pngFileSignature) {
		return Valid
	} else {
		return Invalid
	}
}

func (P PNGParser) GetSize(p []byte) (r Result, t ImageSize) {
	if result := P.DetectType(p); result != Valid {
		return result, ImageSize{}
	}

	if len(p) < 24 {
		return NeedMoreData, ImageSize{}
	}

	i := 8
	for {
		chunkLength := binary.BigEndian.Uint32(p[i:])

		i += 4
		if p[i] == 'I' && p[i+1] == 'H' && p[i+2] == 'D' && p[i+3] == 'R' {
			i += 4
			width := binary.BigEndian.Uint32(p[i:])
			height := binary.BigEndian.Uint32(p[i+4:])
			return Valid, ImageSize{Width: width, Height: height}
		}

		i += 4 + 4 + int(chunkLength) + 4

		if len(p) < i {
			return NeedMoreData, ImageSize{}
		}
	}

}

func init() {
	register(&PNGParser{})
}
