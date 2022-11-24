package parser

// https://datatracker.ietf.org/doc/draft-zern-webp/

type WEBPParser struct{}

func (W WEBPParser) Type() ImageType {
	return WEBP
}

func (W WEBPParser) DetectType(p []byte) (r Result) {
	if len(p) < 12 {
		return NeedMoreData
	}

	if p[0] == 'R' && p[1] == 'I' && p[2] == 'F' && p[3] == 'F' &&
		p[8] == 'W' && p[9] == 'E' && p[10] == 'B' && p[11] == 'P' {
		return Valid
	} else {
		return Invalid
	}
}

func (W WEBPParser) GetSize(p []byte) (r Result, t ImageSize) {
	if result := W.DetectType(p); result != Valid {
		return result, ImageSize{}
	}

	// We need at least 30 bytes for V8PL/VP8X and 24 bytes for V8P
	if len(p) < 30 {
		return NeedMoreData, ImageSize{}
	}

	i := 12

	// VP8 + sync code
	if p[i] == 'V' && p[i+1] == 'P' && p[i+2] == '8' && p[i+3] == ' ' &&
		p[i+11] == '\x9d' && p[i+12] == '\x01' && p[i+13] == '\x2a' {
		i += 14
		width := (uint16(p[i+1])&0x3f)<<8 | uint16(p[i])
		height := (uint16(p[i+3])&0x3f)<<8 | uint16(p[i+2])
		return Valid, ImageSize{Width: uint32(width), Height: uint32(height)}
	}

	// VP8L
	if p[i] == 'V' && p[i+1] == 'P' && p[i+2] == '8' && p[i+3] == 'L' {
		i += 9
		width := 1 + ((uint16(p[i+1])&0x3F)<<8 | uint16(p[i]))
		height := 1 + (uint16(p[i+3])&0xF)<<10 | uint16(p[i+2])<<2 | (uint16(p[i+1])&0xC0)>>6
		return Valid, ImageSize{Width: uint32(width), Height: uint32(height)}
	}

	// VP8X
	if p[i] == 'V' && p[i+1] == 'P' && p[i+2] == '8' && p[i+3] == 'X' {
		i += 12
		width := 1 + (uint32(p[i]) | uint32(p[i+1])<<8 | uint32(p[i+2])<<16)
		height := 1 + (uint32(p[i+3]) | uint32(p[i+4])<<8 | uint32(p[i+5])<<16)
		return Valid, ImageSize{Width: uint32(width), Height: uint32(height)}
	}

	return Invalid, ImageSize{}
}

func init() {
	register(&WEBPParser{})
}
