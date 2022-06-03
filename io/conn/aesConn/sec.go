package aesConn

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"github.com/k773/utils"
	"github.com/k773/utils/io/conn/encryptedConn"
	"strconv"
)

func GenerateCiphers(data, preConfiguredKey []byte, c *encryptedConn.Conn) (e error) {
	if len(data) < 16 || len(data)%16 != 0 || len(data) > (16+32) {
		return utils.NewTemporaryError("GenerateCiphers: wrong data length: " + strconv.Itoa(len(data)))
	}
	var iv = data[:16]
	var key = data[16:]

	if len(preConfiguredKey) != 0 {
		a := sha256.New()
		a.Write(key)
		a.Write(preConfiguredKey)
		key = a.Sum(nil)
	}

	var block cipher.Block
	if block, e = aes.NewCipher(key); e == nil {
		c.Encrypt = cipher.NewCFBEncrypter(block, iv)
		c.Decrypt = cipher.NewCFBDecrypter(block, iv)
	}
	return
}
