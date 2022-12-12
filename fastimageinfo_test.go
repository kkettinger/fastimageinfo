package fastimageinfo

import (
	"github.com/kkettinger/fastimageinfo/parser"
	"testing"
)

func DetectTypeFromFileTesting(filename string, expectedImageType parser.ImageType, t *testing.T) {
	SetChunkSize(1)
	imageType, err := DetectTypeFromFile(filename)
	if err != nil {
		panic(err)
	}

	if imageType != expectedImageType {
		t.Errorf("File %s is expected to be image type %s, but detected type is %s.", filename, expectedImageType, imageType.String())
	}
}

func GetSizeFromFileTesting(filename string, expectedImageSize parser.ImageSize, t *testing.T) {
	SetChunkSize(1)
	imageSize, err := GetSizeFromFile(filename)
	if err != nil {
		panic(err)
	}

	if imageSize.Width != expectedImageSize.Width || imageSize.Height != expectedImageSize.Height {
		t.Errorf("File %s is expected to be of size (%d,%d), but detected size is (%d,%d).",
			filename, expectedImageSize.Width, expectedImageSize.Height,
			imageSize.Width, imageSize.Height)
	}
}

func GetInfoFromFileTesting(filename string, expectedImageType parser.ImageType, expectedImageSize parser.ImageSize, t *testing.T) {
	SetChunkSize(1)
	imageInfo, err := GetInfoFromFile(filename)
	if err != nil {
		panic(err)
	}

	if imageInfo.Type != expectedImageType {
		t.Errorf("File %s is expected to be image type %s, but detected type is %s.",
			filename, expectedImageType, imageInfo.Type.String())
	}

	if imageInfo.Size.Width != expectedImageSize.Width || imageInfo.Size.Height != expectedImageSize.Height {
		t.Errorf("File %s is expected to be of size (%d,%d), but detected size is (%d,%d).",
			filename, expectedImageSize.Width, expectedImageSize.Height,
			imageInfo.Size.Width, imageInfo.Size.Height)
	}
}

type TestCase struct {
	File         string
	expectedType parser.ImageType
	expectedSize parser.ImageSize
}

func Test(t *testing.T) {
	testCases := []TestCase{
		// JPEG
		{File: "testdata/jpeg/example_1.jpg", expectedType: parser.JPEG, expectedSize: parser.ImageSize{Width: 2048, Height: 1536}},
		{File: "testdata/jpeg/example_2.jpg", expectedType: parser.JPEG, expectedSize: parser.ImageSize{Width: 800, Height: 600}},
		{File: "testdata/jpeg/example_3.jpg", expectedType: parser.JPEG, expectedSize: parser.ImageSize{Width: 1, Height: 1}},
		{File: "testdata/jpeg/example_4.jpg", expectedType: parser.JPEG, expectedSize: parser.ImageSize{Width: 275, Height: 297}},

		// PNG
		{File: "testdata/png/example_1.png", expectedType: parser.PNG, expectedSize: parser.ImageSize{Width: 172, Height: 178}},
		{File: "testdata/png/example_2.png", expectedType: parser.PNG, expectedSize: parser.ImageSize{Width: 400, Height: 300}},
		{File: "testdata/png/example_3.png", expectedType: parser.PNG, expectedSize: parser.ImageSize{Width: 386, Height: 395}},

		// GIF
		{File: "testdata/gif/example_1.gif", expectedType: parser.GIF, expectedSize: parser.ImageSize{Width: 250, Height: 297}},
		{File: "testdata/gif/example_2.gif", expectedType: parser.GIF, expectedSize: parser.ImageSize{Width: 217, Height: 217}},

		// BMP
		{File: "testdata/bmp/example_1.bmp", expectedType: parser.BMP, expectedSize: parser.ImageSize{Width: 72, Height: 48}},
		{File: "testdata/bmp/example_2.bmp", expectedType: parser.BMP, expectedSize: parser.ImageSize{Width: 200, Height: 200}},

		// WEBP
		{File: "testdata/webp/example_1.webp", expectedType: parser.WEBP, expectedSize: parser.ImageSize{Width: 550, Height: 368}},
		{File: "testdata/webp/example_2.webp", expectedType: parser.WEBP, expectedSize: parser.ImageSize{Width: 400, Height: 301}},
		{File: "testdata/webp/example_3.webp", expectedType: parser.WEBP, expectedSize: parser.ImageSize{Width: 400, Height: 301}},

		// TIFF
		{File: "testdata/tiff/example_1.tif", expectedType: parser.TIFF, expectedSize: parser.ImageSize{Width: 640, Height: 480}},
		{File: "testdata/tiff/example_2.tif", expectedType: parser.TIFF, expectedSize: parser.ImageSize{Width: 232, Height: 205}},
	}

	for _, testCase := range testCases {
		DetectTypeFromFileTesting(testCase.File, testCase.expectedType, t)
		GetSizeFromFileTesting(testCase.File, testCase.expectedSize, t)
		GetInfoFromFileTesting(testCase.File, testCase.expectedType, testCase.expectedSize, t)
	}
}
