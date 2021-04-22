package main

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func extract(srcImageFile string) []byte {
	inputFile, _ := os.Open(srcImageFile)
	reader := bufio.NewReader(inputFile)
	img, _ := png.Decode(reader)
	inputFile.Close()

	sizeOfMessage := getSizeOfMessageFromImage(img)
	ciphertext := extractFromImage(img, sizeOfMessage, 4)
	return ciphertext
}

func extractFromImage(srcImage image.Image, sizeOfMessage uint32, headerOffset uint32) (ciphertext []byte) {
	var c color.RGBA
	var rgbImage *image.RGBA
	var byteIndex, bitIndex uint32
	var lsb byte

	bounds := srcImage.Bounds()
	width := bounds.Dx()  // width of image
	height := bounds.Dy() // height of image

	rgbImage = image.NewRGBA(image.Rect(0, 0, width, height))              // create image with same dims as inputImage
	draw.Draw(rgbImage, rgbImage.Bounds(), srcImage, bounds.Min, draw.Src) // draw inputImage on the RGBA image

	ciphertext = append(ciphertext, 0) // give slice initial size of 1, with initial value of 0
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c = rgbImage.RGBAAt(x, y)

			// RED
			lsb = getLSB(c.R) // get LSB from the red byte
			ciphertext[byteIndex] = setBitInByte(ciphertext[byteIndex], bitIndex, lsb)
			bitIndex++

			if bitIndex >= 8 { // if the byte has been filled, move on to the next
				bitIndex = 0
				byteIndex++
				if byteIndex >= sizeOfMessage+headerOffset { // if all bytes have been read (including the 4 header bytes)
					return ciphertext[headerOffset : sizeOfMessage+headerOffset]
				}
				ciphertext = append(ciphertext, 0) // add new byte
			}

			// GREEN
			lsb = getLSB(c.G) // get LSB from the red byte
			ciphertext[byteIndex] = setBitInByte(ciphertext[byteIndex], bitIndex, lsb)
			bitIndex++
			if bitIndex >= 8 { // if the byte has been filled, move on to the next
				bitIndex = 0
				byteIndex++
				if byteIndex >= sizeOfMessage+headerOffset { // if all bytes have been read (including the 4 header bytes)
					return ciphertext[headerOffset : sizeOfMessage+headerOffset]
				}
				ciphertext = append(ciphertext, 0) // add new byte
			}

			// BLUE
			lsb = getLSB(c.B) // get LSB from the red byte
			ciphertext[byteIndex] = setBitInByte(ciphertext[byteIndex], bitIndex, lsb)
			bitIndex++
			if bitIndex >= 8 { // if the byte has been filled, move on to the next
				bitIndex = 0
				byteIndex++
				if byteIndex >= sizeOfMessage+headerOffset { // if all bytes have been read (including the 4 header bytes)
					return ciphertext[headerOffset : sizeOfMessage+headerOffset]
				}
				ciphertext = append(ciphertext, 0) // add new byte
			}

		}
	}
	return ciphertext
}

func getSizeOfMessageFromImage(img image.Image) (size uint32) {
	sizeByteArray := extractFromImage(img, 4, 0)
	size = combineFourBytesToInt(sizeByteArray[0], sizeByteArray[1], sizeByteArray[2], sizeByteArray[3])
	return
}

func setBitInByte(targetByte byte, bitIdx uint32, incomingBit byte) byte {
	var outgoingByte byte
	var mask byte = 0x80
	mask = mask >> uint(bitIdx)
	if incomingBit == 0 {
		// set the bit in the byte to 0
		outgoingByte = targetByte & (^mask)
	} else {
		// set the bit in the byte to 1
		outgoingByte = targetByte | mask
	}
	return outgoingByte
}

func getLSB(b byte) byte {
	var mask byte = 1 // 0x01
	if b&mask == 1 {
		return 1
	} else {
		return 0
	}
}

func combineFourBytesToInt(one, two, three, four byte) (i uint32) {
	i = uint32(one)
	i = i << 8
	i = i | uint32(two)
	i = i << 8
	i = i | uint32(three)
	i = i << 8
	i = i | uint32(four)
	return
}
