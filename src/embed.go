/*
Copyright (c) 2021 Mark Ditsworth (@markditsworth)
Copyright (c) 2018 Rafael Passos (@Auyer)
*/
package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func embed(srcImageFile string, dstImageFile string, ciphertext []byte) {
	inputFile, _ := os.Open(srcImageFile)
	reader := bufio.NewReader(inputFile)
	img, _ := png.Decode(reader)
	inputFile.Close()

	writeBuffer := new(bytes.Buffer)
	err := embedToBuffer(writeBuffer, img, ciphertext)
	if err != nil {
		fmt.Printf("Error encoding file %v\n", err)
		return
	}
	outputFile, _ := os.Create(dstImageFile)
	defer outputFile.Close()
	writeBuffer.WriteTo(outputFile)
}

func embedToBuffer(buffer *bytes.Buffer, inputImage image.Image, ciphertext []byte) error {
	var messageLength = uint32(len(ciphertext))
	var color_at_px color.RGBA
	var b byte
	var ok bool
	var rgbImage *image.RGBA // a new RGBA image to create
	bounds := inputImage.Bounds()
	width := bounds.Dx()  // width of image
	height := bounds.Dy() // height of image

	// ensure ciphertext can fit in this image
	err := validateImageSize(width, height, messageLength)
	if err != nil {
		return errors.New("The image is not large enough for the message")
	}

	rgbImage = image.NewRGBA(image.Rect(0, 0, width, height))                // create image with same dims as inputImage
	draw.Draw(rgbImage, rgbImage.Bounds(), inputImage, bounds.Min, draw.Src) // draw inputImage on the RGBA image

	one, two, three, four := splitIntoFourBytes(messageLength)
	ciphertext = append([]byte{four}, ciphertext...)
	ciphertext = append([]byte{three}, ciphertext...)
	ciphertext = append([]byte{two}, ciphertext...)
	ciphertext = append([]byte{one}, ciphertext...)

	ch := make(chan byte, 100) // buffered channel (capacity 100)

	go getNextBitFromMessage(ciphertext, ch) // launch goroutine getting every bit to embed

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			color_at_px = rgbImage.RGBAAt(x, y) // get the color at this pixel (same as the source image since this pixel has yet to be modified)

			// RED
			b, ok = <-ch // collect bit from goroutine
			if !ok {     // if channel is closed, all bits have been encoded, so copy the rest of the photo as normal
				rgbImage.SetRGBA(x, y, color_at_px)
				png.Encode(buffer, rgbImage)
			}
			setLSB(&color_at_px.R, b) // update the LSB of the red value with the current bit from the message

			// GREEN
			b, ok = <-ch
			if !ok { // if channel is closed, all bits have been encoded, so copy the rest of the photo as normal
				rgbImage.SetRGBA(x, y, color_at_px)
				png.Encode(buffer, rgbImage)
				return nil
			}
			setLSB(&color_at_px.G, b)

			// BLUE
			b, ok = <-ch
			if !ok { // if channel is closed, all bits have been encoded, so copy the rest of the photo as normal
				rgbImage.SetRGBA(x, y, color_at_px)
				png.Encode(buffer, rgbImage)
				return nil
			}
			setLSB(&color_at_px.B, b)
			rgbImage.SetRGBA(x, y, color_at_px) // set the color at this pixel with the new encoded value
		}
	}
	err = png.Encode(buffer, rgbImage)
	fmt.Println("err")
	return err
}

func setLSB(b *byte, bit byte) {
	if bit == 1 {
		*b = *b | 1
	} else if bit == 0 {
		var mask byte = 254 //0xFE
		*b = *b & mask
	}
}

func validateImageSize(width int, height int, messageLength uint32) error {
	check := ((uint32(width) * uint32(height) * 3) / 8) - 4
	if check < 4 {
		check = 0
	}
	if check < messageLength + 4 {
		return errors.New("message too large")
	}
	return nil
}

func splitIntoFourBytes(x uint32) (one, two, three, four byte) {
	var mask uint32 = 255         // 0xFF (8 bits)
	one = byte(x >> 24)           // bits 32-25
	two = byte((x >> 16) & mask)  // bits 24-17
	three = byte((x >> 8) & mask) // bits 16-9
	four = byte(x & mask)         // bits 8-1
	return
}

func getNextBitFromMessage(byteArray []byte, ch chan byte) {
	var offsetInBytes, offsetInBitsIntoByte int
	var choiceByte byte
	lenOfMessage := len(byteArray)
	for {
		if offsetInBytes >= lenOfMessage { // after all bytes have been read, close the channel
			close(ch)
			return
		}
		choiceByte = byteArray[offsetInBytes]
		ch <- getBitFromByte(choiceByte, offsetInBitsIntoByte)
		offsetInBitsIntoByte++
		if offsetInBitsIntoByte >= 8 { // after 8 bits, increment to the next byte
			offsetInBitsIntoByte = 0
			offsetInBytes++
		}
	}
}

func getBitFromByte(choiceByte byte, indexInByte int) byte {
	choiceByte = choiceByte << uint(indexInByte)
	var mask byte = 128 // 1000 0000
	bit := mask & choiceByte
	if bit == 128 {
		return 1
	} else {
		return 0
	}
}
