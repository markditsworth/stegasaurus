package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// SHA-256 hash
func hash(key string) [32]byte {
	digest := sha256.Sum256([]byte(key))
	return digest
}

func encrypt(data []byte, digest [32]byte) []byte {
	block, _ := aes.NewCipher(digest[:])
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, digest [32]byte) []byte {
	block, _ := aes.NewCipher(digest[:])
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func encryptFile(in_filename string, digest [32]byte) []byte {
	data, _ := ioutil.ReadFile(in_filename)
	ciphertext := encrypt(data, digest)
	return ciphertext
}

func decryptFile(in_filename string, digest [32]byte) []byte {
	data, _ := ioutil.ReadFile(in_filename)
	plaintext := decrypt(data, digest)
	return plaintext
}

func writeToFile(data []byte, filename string) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(data)
}

func main() {
	var in_filename string
	var out_filename string
	var src_image_file string
	var passkey string
	var direction string
	for i, flag := range os.Args[1:] {
		switch flag {
		case "-i":
			in_filename = os.Args[i+2]
			fmt.Println("Got filename: " + in_filename)
		case "-o":
			out_filename = os.Args[i+2]
			fmt.Println("Output file: " + out_filename)
		case "-p":
			passkey = os.Args[i+2]
			fmt.Println("passkey: " + passkey)
		case "--image":
			src_image_file = os.Args[i+2]
		case "--decrypt":
			direction = "decrypt"
			fmt.Println("decrypting...")
		case "--encrypt":
			direction = "encrypt"
			fmt.Println("encrypting...")
		case "--embed":
			direction = "embed"
			fmt.Println("embedding...")
		case "--extract":
			direction = "extract"
			fmt.Println("extracting...")
		}
	}
	digest := hash(passkey)
	if direction == "encrypt" {
		ciphertext := encryptFile(in_filename, digest)
		writeToFile(ciphertext, out_filename)
	} else if direction == "decrypt" {
		plaintext := decryptFile(in_filename, digest)
		writeToFile(plaintext, out_filename)
	} else if direction == "embed" {
		ciphertext := encryptFile(in_filename, digest)
		embed(src_image_file, out_filename, ciphertext)
	} else if direction == "extract" {
		ciphertext := extract(src_image_file)
		plaintext := decrypt(ciphertext, digest)
		writeToFile(plaintext, out_filename)
	}
}
