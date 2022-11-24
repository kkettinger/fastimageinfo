package parser

import (
	"encoding/binary"
)

// https://www.fileformat.info/format/tiff/egff.htm

type TIFFParser struct{}

type TIFFByteOrder int

const (
	LittleEndian TIFFByteOrder = iota
	BigEndian
)

type TIFFInt int

const (
	Uint16 TIFFInt = iota
	Uint32
)

func (T TIFFParser) Type() ImageType {
	return TIFF
}

func (T TIFFParser) DetectType(p []byte) (r Result) {
	if len(p) < 4 {
		return NeedMoreData
	}

	// Little endian header
	if p[0] == 'I' && p[1] == 'I' && p[2] == '*' && p[3] == '\x00' {
		return Valid
	}

	// Big endian header
	if p[0] == 'M' && p[1] == 'M' && p[2] == '\x00' && p[3] == '*' {
		return Valid
	}

	return Invalid
}

func (T TIFFParser) GetSize(p []byte) (r Result, t ImageSize) {
	if result := T.DetectType(p); result != Valid {
		return result, ImageSize{}
	}

	if len(p) < 18 {
		return NeedMoreData, ImageSize{}
	}

	var byteOrder TIFFByteOrder

	// Detect byte order used inside tiff image
	switch p[0] {
	case 'I':
		byteOrder = LittleEndian
	case 'M':
		byteOrder = BigEndian
	default:
		return Invalid, ImageSize{}
	}

	// Version
	version := TIFFGetInt(byteOrder, Uint16, p[2:])

	if version != 42 {
		return Invalid, ImageSize{}
	}

	// Offset of first IFD
	offsetFirstIFD := TIFFGetInt(byteOrder, Uint32, p[4:])

	if len(p) < offsetFirstIFD+2 {
		return NeedMoreData, ImageSize{}
	}

	i := offsetFirstIFD

	tags := make(map[int]int)

	for {
		// Tag entry count
		tagEntryCount := TIFFGetInt(byteOrder, Uint16, p[i:])

		if len(p) < offsetFirstIFD+2+12*tagEntryCount+4 {
			return NeedMoreData, ImageSize{}
		}

		i += 2

		// Go over tags
		for j := 0; j < tagEntryCount; j++ {
			TIFFParseTag(byteOrder, p[i:i+12], tags)
			i += 12
		}

		// NextIFDOffset
		nextIFDOffset := TIFFGetInt(byteOrder, Uint32, p[i:])
		if nextIFDOffset != 0 {
			i += nextIFDOffset
		}

		break
	}

	// Check if we have collected width and height tag
	// ImageWidth = 256
	width, ok := tags[256]
	if !ok {
		return Invalid, ImageSize{}
	}

	// ImageHeight = 257
	height, ok := tags[257]
	if !ok {
		return Invalid, ImageSize{}
	}

	return Valid, ImageSize{Width: uint32(width), Height: uint32(height)}
}

func TIFFGetInt(byteOrder TIFFByteOrder, intType TIFFInt, p []byte) int {
	switch intType {
	case Uint16:
		{
			if byteOrder == LittleEndian {
				return int(binary.LittleEndian.Uint16(p[0:]))
			} else {
				return int(binary.BigEndian.Uint16(p[0:]))
			}
		}
	case Uint32:
		{
			if byteOrder == LittleEndian {
				return int(binary.LittleEndian.Uint32(p[0:]))
			} else {
				return int(binary.BigEndian.Uint32(p[0:]))
			}
		}
	default:
		panic("Invalid intType given")
	}
}

func TIFFParseTag(byteOrder TIFFByteOrder, p []byte, tags map[int]int) {

	tagIdentifier := TIFFGetInt(byteOrder, Uint16, p[0:])
	dataType := TIFFGetInt(byteOrder, Uint16, p[2:])
	dataCount := TIFFGetInt(byteOrder, Uint32, p[4:])

	// Width & height is stored as single uint16/uint32
	if dataCount != 1 {
		return
	}

	// Data offset
	// We are only interested in tags which are SHORT or LONG
	// 3 = SHORT
	if dataType == 3 {
		dataOffset := TIFFGetInt(byteOrder, Uint16, p[8:])
		tags[tagIdentifier] = dataOffset
	}

	// 4 = LONG
	if dataType == 4 {
		dataOffset := TIFFGetInt(byteOrder, Uint32, p[8:])
		tags[tagIdentifier] = dataOffset
	}
}

func init() {
	register(&TIFFParser{})
}
