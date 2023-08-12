package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

func generatePassword(length int) []byte {
	password := make([]byte, length)
	if _, err := rand.Read(password); err != nil {
		fmt.Println("Error generating password:", err)
		os.Exit(1)
	}
	return password
}

func generateIV() []byte {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		fmt.Println("Error generating IV:", err)
		os.Exit(1)
	}
	return iv
}

func encryptFile(inputFile, outputFile string, keyLength int) ([]byte, []byte) {
	var keySize int
	if keyLength == 128 {
		keySize = 16
	} else if keyLength == 256 {
		keySize = 32
	} else {
		fmt.Println("Invalid key length. Please choose either 128 or 256.")
		os.Exit(1)
	}

	password := generatePassword(keySize)
	iv := generateIV()

	// Read the content of the file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		os.Exit(1)
	}

	// Pad the data to be a multiple of the block size (16 bytes for AES)
	padLen := aes.BlockSize - (len(data) % aes.BlockSize)
	pad := byte(padLen)
	data = append(data, bytes.Repeat([]byte{pad}, padLen)...)

	// Create a new AES cipher block using the key
	block, err := aes.NewCipher(password)
	if err != nil {
		fmt.Println("Error creating AES cipher block:", err)
		os.Exit(1)
	}

	// Create a CBC mode encrypter
	ciphertext := make([]byte, len(data))
	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(ciphertext, data)

	// Write the encrypted data to the output file
	if err := os.WriteFile(outputFile, ciphertext, 0644); err != nil {
		fmt.Println("Error writing to output file:", err)
		os.Exit(1)
	}

	return password, iv
}

func main() {
	var inputFile, outputFile string
	var keyLength int

	fmt.Print("Enter the path of the input file: ")
	fmt.Scan(&inputFile)

	fmt.Print("Enter the path for the output file: ")
	fmt.Scan(&outputFile)

	fmt.Print("Choose encryption key length (128 or 256): ")
	fmt.Scan(&keyLength)

	password, iv := encryptFile(inputFile, outputFile, keyLength)

	fmt.Println("File encrypted successfully.")
	fmt.Println("Password:", hex.EncodeToString(password))
	fmt.Println("IV:", hex.EncodeToString(iv))
}
