package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
)

// config the errors message
// config funcs to help

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Encrypt(text string) (string, error) {
	key := os.Getenv("KEY")
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	iv := []byte(key)[:aes.BlockSize]
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return encode(cipherText), nil
}

func Decrypt(text string) (string, error) {
	key := os.Getenv("KEY")
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	iv := []byte(key)[:aes.BlockSize]
	cipherText := decode(text)
	cfb := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func Response(res any, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}
