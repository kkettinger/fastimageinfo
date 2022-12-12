package main

import (
	"fmt"
	"github.com/kkettinger/fastimageinfo"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: " + os.Args[0] + " image.jpg")
		os.Exit(1)
		return
	}

	imageInfo, err := fastimageinfo.GetInfoFromFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	fmt.Printf("File:\t%s\n", os.Args[1])
	fmt.Printf("Type:\t%s\n", imageInfo.Type.String())
	fmt.Printf("Size:\t%d x %d\n", imageInfo.Size.Width, imageInfo.Size.Height)
	fmt.Printf("Mime:\t%s\n", imageInfo.Type.ToMimetype())
}
