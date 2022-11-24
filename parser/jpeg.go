package parser

import (
	"bytes"
)

// Information about the jpeg structure can be found here:
// https://www.ccoderun.ca/programming/2017-01-31_jpeg/
// https://www.w3.org/Graphics/JPEG/itu-t81.pdf

type JPEGParser struct{}

func (J JPEGParser) Type() ImageType {
	return JPEG
}

func (J JPEGParser) DetectType(p []byte) (r Result) {
	// SOI
	startOfImage := []byte{'\xff', '\xd8'}

	if len(p) < len(startOfImage) {
		return NeedMoreData
	}

	if !bytes.Equal(p[:len(startOfImage)], startOfImage) {
		return Invalid
	}

	return Valid
}

// Credits go to https://web.archive.org/web/20130305080105/http://www.64lines.com/jpeg-width-height
func (J JPEGParser) GetSize(p []byte) (r Result, t ImageSize) {
	if result := J.DetectType(p); result != Valid {
		return result, ImageSize{}
	}

	if len(p) < 6 {
		return NeedMoreData, ImageSize{}
	}

	/*jfifStart := 6
	jfifIdentifier := []byte{'J', 'F', 'I', 'F', '\x00'}

	if !enoughData(p, jfifStart+len(jfifIdentifier)) {
		return NeedMoreData, ImageSize{}
	}

	// Check for jfif identifier
	if !bytes.Equal(p[jfifStart:jfifStart+len(jfifIdentifier)], jfifIdentifier) {
		fmt.Println("No jfif tag found")
		return Invalid, ImageSize{}
	}*/

	var i uint32 = 4

	// List of valid SOFs
	// SOF0  = 0xC0 Baseline DCT
	// SOF1  = 0xC1 Extended sequential DCT, Huffman coding
	// SOF2  = 0xC2 Progressive DCT, Huffman coding
	// SOF3  = 0xC3 Lossless (sequential), Huffman coding
	// SOF9  = 0xC9 Extended sequential DCT, arithmetic coding
	// SOF10 = 0xCA Progressive DCT, arithmetic coding
	// SOF11 = 0xCB Lossless (sequential), arithmetic coding
	validSOFs := []byte{'\xc0', '\xc1', '\xc2', '\xc3', '\xc9', '\xca', '\xcb'}

	// Retrieve the block length of the first block since the first block will not contain the size of file
	blockLength := uint32(p[i])*256 + uint32(p[i+1])
	for {
		// Increase the file index to get to the next block
		i += blockLength

		// Check if we have enough data
		if len(p) < int(i)+9 {
			return NeedMoreData, ImageSize{}
		}

		// Check that we are truly at the start of another block
		if p[i] != '\xff' {
			return Invalid, ImageSize{}
		}

		// 0xFF<SOFN> is the "Start of frame" marker which contains the file size
		if bytes.Contains(validSOFs, []byte{p[i+1]}) {
			// The structure of the 0xFFC0 block is quite simple
			// [0xFF<SOFN>][ushort length][uchar precision][ushort x][ushort y]
			height := uint32(p[i+5])*256 + uint32(p[i+6])
			width := uint32(p[i+7])*256 + uint32(p[i+8])
			return Valid, ImageSize{Width: width, Height: height}
		} else {
			// Skip the block marker
			i += 2

			// Go to the next block
			blockLength = uint32(p[i])*256 + uint32(p[i+1])
		}
	}
}

func init() {
	register(&JPEGParser{})
}
