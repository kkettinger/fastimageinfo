# fastimageinfo ![build](https://github.com/hberg539/fastimageinfo/actions/workflows/go.yml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/hberg539/fastimageinfo.svg)](https://pkg.go.dev/github.com/hberg539/fastimageinfo)
fastimageinfo is a go module for fetching image type, width and height without reading the whole file.

## Install
```bash
$ go get github.com/kkettinger/fastimageinfo
```

## How it works
fastimageinfo reads multiple byte chunks until the image type, width and height could be detected.

## How to use

- To get type, width and height, use `GetInfo()`, `GetInfoFromReader()`, `GetInfoFromFile()`
- To only detect the type, use `DetectType()`, `DetectTypeFromReader()`, `DetectTypeFromFile()`
- The chunk size used by the `*FromReader()` and `*FromFile()` functions can be set with `SetChunkSize(byte)`

###  Example: Read from file

```go
imageInfo, err := fastimageinfo.GetInfoFromFile("testdata/jpeg/example_1.jpg")
if err != nil {
    panic(err)
}

fmt.Println(imageInfo.Type)
fmt.Println(imageInfo.Size)
fmt.Println(imageInfo.Type.ToMimetype())
```

Output:
```
JPEG      
{1050 700}
image/jpeg
```

### Example: Read from file handle
Here an example by passing an open file handle:

```go
reader, err := os.Open("testdata/jpeg/example_1.jpg")
defer reader.Close()
if err != nil {
    panic(err)
}

imageInfo, bytesRead, err := fastimageinfo.GetInfoFromReader(reader)
if err != nil {
    panic(err)
}

fmt.Println(imageInfo.Type)
fmt.Println(imageInfo.Size)
fmt.Println(bytesRead)
```

You can control how big the chunks are with `fastimageinfo.SetChunkSize(bytes)`.

### Example: Read from http response
Here an example with a 34MB jpg file:

```go
resp, err := http.Get("https://upload.wikimedia.org/wikipedia/commons/5/5e/M104_ngc4594_sombrero_galaxy_hi-res.jpg")
defer resp.Body.Close()
if err != nil {
    panic(err)
}

fastimageinfo.SetChunkSize(1)
imageInfo, bytesRead, err := fastimageinfo.GetInfoFromReader(resp.Body)
if err != nil {
    panic(err)
}

fmt.Println(imageInfo.Type)
fmt.Println(imageInfo.Size)
fmt.Println(bytesRead)
```

Output:
```
JPEG
{11472 6429}
232
```

Only 232 bytes were needed to find out the image type and dimension for this image.


### Example: Read from byte array
You can also pass a byte stream and the `GetInfo()` method will tell you if it needs more data to make a decision.

```go
// First bytes of a 1x1px jpg
data := []byte{
    '\xFF', '\xD8', '\xFF', '\xE0', '\x00', '\x10', '\x4A', '\x46', '\x49', '\x46', '\x00', '\x01', '\x01', '\x01', '\x00', '\x60',
    '\x00', '\x60', '\x00', '\x00', '\xFF', '\xE1', '\x00', '\x5A', '\x45', '\x78', '\x69', '\x66', '\x00', '\x00', '\x4D', '\x4D',
    '\x00', '\x2A', '\x00', '\x00', '\x00', '\x08', '\x00', '\x05', '\x03', '\x01', '\x00', '\x05', '\x00', '\x00', '\x00', '\x01',
    '\x00', '\x00', '\x00', '\x4A', '\x03', '\x03', '\x00', '\x01', '\x00', '\x00', '\x00', '\x01', '\x00', '\x00', '\x00', '\x00',
    '\x51', '\x10', '\x00', '\x01', '\x00', '\x00', '\x00', '\x01', '\x01', '\x00', '\x00', '\x00', '\x51', '\x11', '\x00', '\x04',
    '\x00', '\x00', '\x00', '\x01', '\x00', '\x00', '\x0E', '\xC3', '\x51', '\x12', '\x00', '\x04', '\x00', '\x00', '\x00', '\x01',
    '\x00', '\x00', '\x0E', '\xC3', '\x00', '\x00', '\x00', '\x00', '\x00', '\x01', '\x86', '\xA0', '\x00', '\x00', '\xB1', '\x8F',
    '\xFF', '\xDB', '\x00', '\x43', '\x00', '\x02', '\x01', '\x01', '\x02', '\x01', '\x01', '\x02', '\x02', '\x02', '\x02', '\x02',
    '\x02', '\x02', '\x02', '\x03', '\x05', '\x03', '\x03', '\x03', '\x03', '\x03', '\x06', '\x04', '\x04', '\x03', '\x05', '\x07',
    '\x06', '\x07', '\x07', '\x07', '\x06', '\x07', '\x07', '\x08', '\x09', '\x0B', '\x09', '\x08', '\x08', '\x0A', '\x08', '\x07',
    '\x07', '\x0A', '\x0D', '\x0A', '\x0A', '\x0B', '\x0C', '\x0C', '\x0C', '\x0C', '\x07', '\x09', '\x0E', '\x0F', '\x0D', '\x0C',
    '\x0E', '\x0B', '\x0C', '\x0C', '\x0C', '\xFF', '\xDB', '\x00', '\x43', '\x01', '\x02', '\x02', '\x02', '\x03', '\x03', '\x03',
    '\x06', '\x03', '\x03', '\x06', '\x0C', '\x08', '\x07', '\x08', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C',
    '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C',
    '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C',
    '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\x0C', '\xFF', '\xC0', '\x00', '\x11', '\x08', '\x00',
    '\x01', '\x00', '\x01'}

result, imageInfo, err := fastimageinfo.GetInfo(data)
if err != nil {
    panic(err)
}

// Check if the parser has enough data to make a decision and/or extract image dimensions
if result == fastimageinfo.NeedMoreData {
    // .... feed more data
    panic("Not enough data to make a decision")
}

// Invalid means, type and size could not be determined
if result == fastimageinfo.Invalid {
    panic("Image type and/or size could not be determined")
}

// Image type and size could be determined
if result == fastimageinfo.Valid {
    fmt.Println(imageInfo.Type)
    fmt.Println(imageInfo.Size)
}
```

## Supported image types

- JPEG
- PNG
- BMP
- GIF
- WEBP
- TIFF
