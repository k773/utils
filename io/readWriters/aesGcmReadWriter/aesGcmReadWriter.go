package aesGcmReadWriter

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/k773/utils/io/readWriters/packetReadWriter"
	"github.com/pkg/errors"
)

/*
	aesGcmReadWriter packet structure: [ciphertext..., nonce(12)]
*/

type ReadWriter struct {
	PacketReadWriter *packetReadWriter.ReadWriter
	BlockCipher      cipher.AEAD

	// optional
	getNonce func() []byte
}

func NewReadWriter(packetReadWriter *packetReadWriter.ReadWriter, blockCipher cipher.AEAD, getNonce func() []byte) *ReadWriter {
	if getNonce == nil {
		getNonce = func() []byte {
			var nonce = make([]byte, 12)
			if _, e := rand.Read(nonce); e != nil {
				panic(e)
			}
			return nonce
		}
	}
	return &ReadWriter{PacketReadWriter: packetReadWriter, BlockCipher: blockCipher, getNonce: getNonce}
}

var ErrorIncorrectEncryptedPacketSize = errors.New("incorrect encrypted packet size")

func (rw *ReadWriter) ReadPacket() (packet *bytes.Buffer, err error) {
	ciphertext, err := rw.PacketReadWriter.ReadPacket(nil)
	if err == nil {
		if ciphertext.Len() < 12 {
			return nil, ErrorIncorrectEncryptedPacketSize
		}

		data := ciphertext.Bytes()
		nonceStart := ciphertext.Len() - 12
		data, err = rw.BlockCipher.Open(data[:0], data[nonceStart:], data[:nonceStart], nil)
		packet = bytes.NewBuffer(data)
	}
	return
}

func (rw *ReadWriter) WritePacket(packet *bytes.Buffer) (err error) {
	var dst = make([]byte, 16+12+packet.Len())
	var nonceStart = len(dst) - 12
	copy(dst[nonceStart:], rw.getNonce())

	rw.BlockCipher.Seal(dst[:0], dst[nonceStart:], packet.Bytes(), nil)
	return rw.PacketReadWriter.WritePacket(bytes.NewBuffer(dst))
}
