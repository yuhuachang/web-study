package main

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
)

// https://en.wikipedia.org/wiki/ICO_(file_format)
func createIcoImage(icoSize int) []byte {

	// Create ICO header
	icoHeaderSize := 6
	icoHeader := make([]byte, icoHeaderSize)
	binary.LittleEndian.PutUint16(icoHeader[0:2], 0) // Reserved. Must always be 0.
	binary.LittleEndian.PutUint16(icoHeader[2:4], 1) // Specifies image type: 1 for icon (.ICO) image, 2 for cursor (.CUR) image. Other values are invalid.
	binary.LittleEndian.PutUint16(icoHeader[4:6], 1) // Specifies number of images in the file.

	// // Create ICO image directory entry
	icoImageDirectoryEntrySize := 16
	icoImageDirectoryEntry := make([]byte, icoImageDirectoryEntrySize)
	icoImageDirectoryEntry[0] = uint8(icoSize)                                                                     // Specifies image width in pixels. Can be any number between 0 and 255. Value 0 means image width is 256 pixels.
	icoImageDirectoryEntry[1] = uint8(icoSize)                                                                     // Specifies image height in pixels. Can be any number between 0 and 255. Value 0 means image height is 256 pixels.
	icoImageDirectoryEntry[2] = 0                                                                                  // Specifies number of colors in the color palette. Should be 0 if the image does not use a color palette.
	icoImageDirectoryEntry[3] = 0                                                                                  // Reserved. Should be 0.
	binary.LittleEndian.PutUint16(icoImageDirectoryEntry[4:6], 1)                                                  // In ICO format: Specifies color planes. Should be 0 or 1.  In CUR format: Specifies the horizontal coordinates of the hotspot in number of pixels from the left.
	binary.LittleEndian.PutUint16(icoImageDirectoryEntry[6:8], 32)                                                 // In ICO format: Specifies bits per pixel. Typical values are 1, 4, 8, 16, 24 and 32. In CUR format: Specifies the vertical coordinates of the hotspot in number of pixels from the top.
	binary.LittleEndian.PutUint16(icoImageDirectoryEntry[8:12], uint16(icoSize*icoSize))                           // Specifies the size of the image's data in bytes
	binary.LittleEndian.PutUint16(icoImageDirectoryEntry[12:16], uint16(icoHeaderSize+icoImageDirectoryEntrySize)) // Specifies the offset of BMP or PNG data from the beginning of the ICO/CUR file

	// Draw image
	img := image.NewRGBA(image.Rect(0, 0, icoSize, icoSize))
	for x := 0; x < icoSize; x++ {
		for y := 0; y < icoSize; y++ {
			r := uint8(0xff / icoSize * x)
			g := uint8(0xff / icoSize * y)
			img.Set(x, y, color.RGBA{r, g, 0xff, 0xff})
		}
	}

	var imageBuffer bytes.Buffer
	png.Encode(&imageBuffer, img)

	// Create ICO file
	icoFile := []byte{}
	icoFile = append(icoFile, icoHeader...)
	icoFile = append(icoFile, icoImageDirectoryEntry...)
	icoFile = append(icoFile, imageBuffer.Bytes()...)

	return icoFile
}

func main() {

	icoFile := createIcoImage(32)

	// Write ICO file
	ioutil.WriteFile("favicon.ico", icoFile, 0644)
}
