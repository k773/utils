package chacha20_poly1305_conn

import (
	"encoding/binary"
	"io"
)

func (c *Conn) sendPacketNoLock(header byte, data, cipherAdditionalData []byte) (e error) {
	switch header {
	case headerRekey:
		_, e = c.Underlying.Write([]byte{header})
	default:
		var nonce []byte
		if c.cipher != nil {
			var nonceSize = c.cipher.NonceSize()
			nonce = generateNonce(nonceSize)

			var buf = make([]byte, 8+nonceSize+8+len(data)+c.cipher.Overhead())
			// Set nonceSize && nonce
			binary.LittleEndian.PutUint64(buf, uint64(nonceSize))
			copy(buf[8:], nonce)
			// Set encrypted data
			data = c.cipher.Seal(data[:0], nonce, data, cipherAdditionalData)
			binary.LittleEndian.PutUint64(buf[8+nonceSize:], uint64(len(data)))
			writtenEncData := copy(buf[8+nonceSize+8:], data)
			// Stripping the buffer to the actual size
			buf = buf[:8+nonceSize+8+writtenEncData]
			data = buf
		} else {
			// data: bytes(len(data)) + data
			var buf = make([]byte, 8+len(data))
			binary.LittleEndian.PutUint64(buf, uint64(len(data)))
			copy(buf[8:], data)
			data = buf
		}
		if e == nil {
			if _, e = c.Underlying.Write([]byte{header}); e == nil {
				_, e = c.Underlying.Write(data)
			}
		}
	}

	return
}

func (c *Conn) mustReadNoLockNoDecryption(n int) (data []byte, e error) {
	data = make([]byte, n)
	_, e = io.ReadFull(c.Underlying, data)
	return
}

func (c *Conn) readPacketNoLock(cipherAdditionalData []byte) (header byte, data []byte, e error) {
	var packetLength = 0

	var onPrePacketReceive = func(c *Conn, size int) (e error) {
		if c.cipher != nil {
			packetLength += size
			c.bytesSendReceivedEnc += size
			if packetLength > MaxPacketLength || c.bytesSendReceivedEnc >= RejectAfterByte {
				e = ErrorTooLargePacket
			}
		}
		return
	}
	var decrypt = func(c *Conn, nonce, data, cipherAdditionalData []byte) (res []byte, e error) {
		if nonce == nil {
			return data, nil
		}
		res, e = c.cipher.Open(data[:0], nonce, data, cipherAdditionalData)
		return
	}
	var readSinglePacket = func(c *Conn, nonce, cipherAdditionalData []byte) (data []byte, e error) {
		var packetSizeBytes []byte
		if e = onPrePacketReceive(c, 8); e == nil {
			if packetSizeBytes, e = c.mustReadNoLockNoDecryption(8); e == nil {
				var packetSize = int(binary.LittleEndian.Uint64(packetSizeBytes))
				if e = onPrePacketReceive(c, packetSize); e == nil {
					if data, e = c.mustReadNoLockNoDecryption(packetSize); e == nil {
						data, e = decrypt(c, nonce, data, cipherAdditionalData)
					}
				}
			}
		}
		return
	}

	headerBytes, e := c.mustReadNoLockNoDecryption(1)
	if e == nil {
		switch header = headerBytes[0]; header {
		case headerRekey:
			break
		case headerData:
			if c.cipher != nil {
				var nonce []byte
				if nonce, e = readSinglePacket(c, nil, nil); e == nil {
					data, e = readSinglePacket(c, nonce, cipherAdditionalData)
				}
			} else {
				data, e = readSinglePacket(c, nil, nil)
			}
		default:
			e = ErrorUnknownHeaderReceived
		}
	}
	return
}

func (c *Conn) readDataPacketNoLock(cipherAdditionalData []byte) (data []byte, e error) {
	header, data, e := c.readPacketNoLock(cipherAdditionalData)
	if header != headerData {
		e = ErrorWrongHeaderReceived
	}
	return
}
