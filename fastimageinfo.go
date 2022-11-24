package fastimageinfo

import (
	"bytes"
	"github.com/hberg539/fastimageinfo/parser"
	"io"
	"os"
)

type ImageInfo struct {
	Type parser.ImageType
	Size parser.ImageSize
}

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

var chunkSize = 128

func SetChunkSize(s int) {
	chunkSize = s
}

func DetectType(p []byte) (Result, parser.ImageType, error) {
	notEnoughData := false

	for imageType, imageParser := range parser.ImageParsers {
		result := imageParser.DetectType(p)

		if result == parser.NeedMoreData {
			notEnoughData = true
		}

		if result == parser.Valid {
			return Valid, imageType, nil
		}
	}

	if notEnoughData {
		return NeedMoreData, parser.UnknownType, nil
	}

	return Invalid, parser.UnknownType, nil
}

func GetSize(p []byte) (Result, parser.ImageSize, error) {
	result, imageType, err := DetectType(p)
	if err != nil || result != Valid {
		return result, parser.ImageSize{}, err
	}

	resultParser, imageSize := parser.ImageParsers[imageType].GetSize(p)

	if resultParser == parser.NeedMoreData {
		return NeedMoreData, parser.ImageSize{}, nil
	}

	if resultParser == parser.Valid {
		return Valid, imageSize, nil
	}

	return Invalid, parser.ImageSize{}, nil
}

func GetInfo(p []byte) (Result, ImageInfo, error) {
	result, imageType, err := DetectType(p)
	if err != nil || result != Valid {
		return result, ImageInfo{}, err
	}

	result, imageSize, err := GetSize(p)
	if err != nil || result != Valid {
		return result, ImageInfo{}, err
	}

	imageInfo := ImageInfo{
		Type: imageType,
		Size: imageSize,
	}

	return Valid, imageInfo, nil
}

func DetectTypeFromReader(r io.Reader) (parser.ImageType, int, error) {
	buf := bytes.Buffer{}

	for {
		chunk := make([]byte, chunkSize)

		count, err := r.Read(chunk)
		if err != nil {
			return parser.UnknownType, 0, err
		}

		buf.Write(chunk[:count])

		result, imageType, err := DetectType(buf.Bytes())
		if err != nil || result == Invalid || result == Valid {
			return imageType, len(buf.Bytes()), err
		}

		if result == NeedMoreData {
			continue
		}
	}
}

func GetSizeFromReader(r io.Reader) (parser.ImageSize, int, error) {
	buf := bytes.Buffer{}

	for {
		chunk := make([]byte, chunkSize)

		count, err := r.Read(chunk)
		if err != nil {
			return parser.ImageSize{}, 0, err
		}

		buf.Write(chunk[:count])

		result, imageSize, err := GetSize(buf.Bytes())
		if err != nil || result == Invalid || result == Valid {
			return imageSize, len(buf.Bytes()), err
		}

		if result == NeedMoreData {
			continue
		}
	}
}

func GetInfoFromReader(r io.Reader) (ImageInfo, int, error) {
	buf := bytes.Buffer{}
	for {
		chunk := make([]byte, chunkSize)

		count, err := r.Read(chunk)
		if err != nil {
			return ImageInfo{}, 0, err
		}

		buf.Write(chunk[:count])

		result, imageInfo, err := GetInfo(buf.Bytes())
		if err != nil || result == Invalid || result == Valid {
			return imageInfo, len(buf.Bytes()), err
		}

		if result == NeedMoreData {
			continue
		}
	}
}

func DetectTypeFromFile(filepath string) (parser.ImageType, error) {
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return parser.UnknownType, err
	}

	imageType, _, err := DetectTypeFromReader(f)
	return imageType, err
}

func GetSizeFromFile(filepath string) (parser.ImageSize, error) {
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return parser.ImageSize{}, err
	}

	imageSize, _, err := GetSizeFromReader(f)
	return imageSize, err
}

func GetInfoFromFile(filepath string) (ImageInfo, error) {
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return ImageInfo{}, err
	}

	imageInfo, _, err := GetInfoFromReader(f)
	return imageInfo, err
}
