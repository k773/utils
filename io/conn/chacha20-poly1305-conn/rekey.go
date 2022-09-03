package chacha20_poly1305_conn

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"github.com/k773/utils"
	"golang.org/x/crypto/chacha20poly1305"
)

func (c *Conn) rekey(wLocked, rLocked bool) (e error) {
	if !wLocked {
		c.w.Lock()
		defer c.w.Unlock()
	}
	if !rLocked {
		c.r.Lock()
		defer c.r.Unlock()
	}

	c.onPreRekey(c)

	if c.IsServer {
		e = c.rekeyServersideNoLock()
	} else {
		e = c.rekeyClientSideNoLock()
	}

	c.onPostRekey(c, e)
	return
}

/*
rekeyClientSide initiates a client-side rekey procedure:
1) receives 256 bits from the server
2) computes a new session key: sha256(key + received data)
3) receives a message from the server -> decrypt -> reverse bytes -> encrypt -> send back
4) receives the verification result from the server

If rekeyClientSide has failed, then the cipher is not set.
*/
func (c *Conn) rekeyClientSideNoLock() (e error) {
	// Dropping the expired cipher
	c.cipher = nil
	c.bytesSendReceivedEnc = 0

	// 1)
	randomData, e := c.readDataPacketNoLock(nil)
	if e == nil {
		// 2)
		if e = c.generateSessionCipher(randomData); e == nil {
			// 3&4)
			e = c.verifyConnectionClientside()
		}
	}

	if e != nil {
		c.cipher = nil
		c.bytesSendReceivedEnc = 0
	}
	return
}

/*
rekeyServerside initiates a rekey procedure:
1) generates a random 256 bits and sends them to the client
2) computes a new session key: sha256(key + generated data)
3) verifies the connection: generates and sends 256 random bits to the client (by the encrypted channel) and expects to see them in reverse (after the decryption)
4) sends verification result to the client

If rekeyServerside has failed, then the cipher is not set.
*/
func (c *Conn) rekeyServersideNoLock() (e error) {
	// Dropping the expired cipher
	c.cipher = nil
	c.bytesSendReceivedEnc = 0

	// 1)
	var randomData = make([]byte, 64)
	if _, e = rand.Read(randomData); e == nil {
		if e = c.sendPacketNoLock(headerData, randomData, nil); e == nil {
			// 2)
			if e = c.generateSessionCipher(randomData); e == nil {
				// 3&4)
				e = c.verifyConnectionServerside()
			}
		}
	}

	if e != nil {
		c.cipher = nil
		c.bytesSendReceivedEnc = 0
	}
	return
}

func (c *Conn) generateSessionCipher(randomData []byte) (e error) {
	var hash = sha256.New()
	if _, e = hash.Write(c.baseKey); e == nil {
		if _, e = hash.Write(randomData); e == nil {
			c.cipher, e = chacha20poly1305.New(hash.Sum(nil))
		}
	}
	return
}

func (c *Conn) verifyConnectionServerside() (e error) {
	var randomData = make([]byte, 64)
	if _, e = rand.Read(randomData); e == nil {
		if e = c.sendPacketNoLock(headerData, randomData, []byte("verify_data")); e == nil {
			var received []byte
			if received, e = c.readDataPacketNoLock([]byte("verify_data_response")); e == nil {
				var verificationResultByte byte = 1
				if !bytes.Equal(received, utils.Reverse(randomData)) {
					e = ErrorVerificationFailed
					verificationResultByte = 0
				}
				_ = c.sendPacketNoLock(headerData, []byte{verificationResultByte}, []byte("verify_data_result"))
			}
		}
	}
	return
}

func (c *Conn) verifyConnectionClientside() (e error) {
	data, e := c.readDataPacketNoLock([]byte("verify_data"))
	if e == nil {
		if e = c.sendPacketNoLock(headerData, utils.Reverse(data), []byte("verify_data_response")); e == nil {
			if data, e = c.readDataPacketNoLock([]byte("verify_data_result")); e == nil {
				if len(data) != 1 || data[0] != 1 {
					e = ErrorVerificationFailed
				}
			}
		}
	}
	return
}
