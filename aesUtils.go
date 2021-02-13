package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func Encrypt(key, data []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	cipherText := make([]byte, aes.BlockSize+len(data))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], data)

	return cipherText
}

func Decrypt(key, data, iv []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	return data
}

func DecryptBtB(hexKey string, bytes []byte) []byte {
	return Decrypt(H2b(hexKey), bytes[aes.BlockSize:], bytes[:aes.BlockSize])
}

func EncryptBtB(hexKey string, data []byte) []byte {
	return Encrypt(H2b(hexKey), data)
}

func EncryptBtBSafe(aesKey string, data []byte) []byte {
	randomBytes := make([]byte, 16)
	_, _ = rand.Read(randomBytes)
	return EncryptBtB(aesKey, append(randomBytes, data...))
}

func DecryptBtBSafe(aesKey string, data []byte) []byte {
	return DecryptBtB(aesKey, data)[16:]
}

func DecryptHtB(hexKey, hexData string) []byte {
	data := H2b(hexData)
	//fmt.Println(hexData, data)
	return Decrypt(H2b(hexKey), data[aes.BlockSize:], data[:aes.BlockSize])
}

func EncryptBtH(hexKey string, data []byte) string {
	return hex.EncodeToString(Encrypt(H2b(hexKey), data))
}

func DecryptHtS(hexKey, hexStr string) string {
	cipherText := H2b(hexStr)
	return string(Decrypt(H2b(hexKey), cipherText[aes.BlockSize:], cipherText[:aes.BlockSize]))
}

func DecryptB64tB(hexKey, b64 string) []byte {
	cipherText, _ := base64.StdEncoding.DecodeString(b64)
	return Decrypt(H2b(hexKey), cipherText[aes.BlockSize:], cipherText[:aes.BlockSize])
}

func EncryptBtB64(hexKey string, data []byte) string {
	return base64.StdEncoding.EncodeToString(Encrypt(H2b(hexKey), data))
}

func DecryptB64tBSafe(hexKey, b64 string) []byte {
	return DecryptB64tB(hexKey, b64)[16:]
}

func EncryptBtB64Safe(hexKey string, data []byte) string {
	iv := make([]byte, 16)
	rand.Read(iv)
	return EncryptBtB64(hexKey, append(iv, data...))
}
