/*
	The package provides a chacha20 implementation of io.ReadWriter
*/

package chacha20UnauthenticatedReadWriter

import (
	"golang.org/x/crypto/chacha20"
	"io"
)

type ReadWriter struct {
	Back               io.ReadWriter
	BlockCipherEncrypt *chacha20.Cipher
	BlockCipherDecrypt *chacha20.Cipher
}

func NewReadWriter(back io.ReadWriter, blockCipherEncrypt, blockCipherDecrypt *chacha20.Cipher) *ReadWriter {
	return &ReadWriter{Back: back, BlockCipherEncrypt: blockCipherEncrypt, BlockCipherDecrypt: blockCipherDecrypt}
}

func (rw *ReadWriter) Read(p []byte) (n int, err error) {
	n, err = rw.Back.Read(p)
	if rw.BlockCipherDecrypt != nil && n != 0 {
		rw.BlockCipherDecrypt.XORKeyStream(p[:n], p[:n])
	}
	return
}

func (rw *ReadWriter) Write(p []byte) (n int, err error) {
	if rw.BlockCipherEncrypt != nil {
		rw.BlockCipherEncrypt.XORKeyStream(p, p)
	}
	return rw.Back.Write(p)
}
