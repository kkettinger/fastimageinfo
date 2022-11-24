package parser

type Result int

const (
	UnknownResult Result = iota
	NeedMoreData
	Valid
	Invalid
)

func (p Result) String() string {
	switch p {
	case NeedMoreData:
		return "NeedMoreData"
	case Valid:
		return "Valid"
	case Invalid:
		return "Invalid"
	case UnknownResult:
		return "UnknownResult"
	default:
		return "UnknownResult"
	}
}

type ImageType int

const (
	UnknownType ImageType = iota
	JPEG
	PNG
	BMP
	GIF
	WEBP
	TIFF
)

func (t ImageType) String() string {
	switch t {
	case JPEG:
		return "JPEG"
	case PNG:
		return "PNG"
	case BMP:
		return "BMP"
	case GIF:
		return "GIF"
	case WEBP:
		return "WEBP"
	case TIFF:
		return "TIFF"
	case UnknownType:
		return "UnknownType"
	default:
		return "UnknownType"
	}
}

func (t ImageType) ToMimetype() string {
	switch t {
	case JPEG:
		return "image/jpeg"
	case PNG:
		return "image/png"
	case BMP:
		return "image/bmp"
	case GIF:
		return "image/gif"
	case WEBP:
		return "image/webp"
	case TIFF:
		return "image/tiff"
	case UnknownType:
		return "application/octet-stream"
	default:
		return "application/octet-stream"
	}
}

type ImageSize struct {
	Width  uint32
	Height uint32
}

type ImageParser interface {
	Type() ImageType
	DetectType(p []byte) (r Result)
	GetSize(p []byte) (r Result, t ImageSize)
}

var ImageParsers = make(map[ImageType]ImageParser)

func register(imageParser ImageParser) {
	ImageParsers[imageParser.Type()] = imageParser
}
