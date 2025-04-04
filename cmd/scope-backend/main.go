package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"scope/internal/utils"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello World From Scope Backend")
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	ciphertext, err := utils.Encrypt([]byte("Hello World From Scope Backend"), os.Getenv("AES256_KEY"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Ciphertext:", hex.EncodeToString(ciphertext))
	plaintext, err := utils.Decrypt(ciphertext, os.Getenv("AES256_KEY"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Plaintext:", string(plaintext))
}
