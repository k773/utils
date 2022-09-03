package chacha20_poly1305_conn

import (
	"crypto/cipher"
	"errors"
	"net"
	"sync"
)

/*
	Rekey constants
	TODO: synchronize between client and server
*/

const (
	RekeyAfterByte  = 1 << 30            // Rekey every 1gb of both send and received data
	RejectAfterByte = 2 * RekeyAfterByte // The client will drop all packets after RejectAfterByte
	MaxPacketLength = RejectAfterByte
)

var (
	ErrorTooLargePacket        = errors.New("too large packet")
	ErrorTooHighCounterState   = errors.New("too high counter state")
	ErrorRecursiveRekey        = errors.New("received two rekey requests in a row")
	ErrorUnknownHeaderReceived = errors.New("unknown header received")
	ErrorWrongHeaderReceived   = errors.New("wrong header received")
	ErrorVerificationFailed    = errors.New("verification failed")
)

/*
	Headers
*/

const (
	headerRekey byte = 0
	headerData  byte = 1
)

type Conn struct {
	Underlying net.Conn
	// Defines, what rekey procedure should be used by this Conn instance
	IsServer bool

	bytesSendReceivedEnc int

	cipher      cipher.AEAD
	baseKey     []byte
	onPreRekey  func(c *Conn)
	onPostRekey func(c *Conn, e error)

	r, w sync.Mutex
}

/*
	Upgrade
*/

func Upgrade(underlying net.Conn, isServer bool, baseKey []byte, onPreRekey func(c *Conn), onPostRekey func(c *Conn, e error)) (c *Conn, e error) {
	if onPreRekey == nil {
		onPreRekey = func(c *Conn) {}
	}
	if onPostRekey == nil {
		onPostRekey = func(c *Conn, e error) {}
	}
	c = &Conn{
		Underlying:  underlying,
		IsServer:    isServer,
		baseKey:     baseKey,
		onPreRekey:  onPreRekey,
		onPostRekey: onPostRekey,
	}
	return c, c.rekey(false, false)
}

/*
	Read && write
*/

func (c *Conn) SendPacket(data, cipherAdditionalData []byte) (e error) {
	c.w.Lock()
	defer c.w.Unlock()

	var doSend bool

	if c.bytesSendReceivedEnc+len(data) >= RejectAfterByte {
		if len(data) <= RejectAfterByte {
			if e = c.rekey(true, false); e == nil {
				doSend = true
			}
		} else {
			e = ErrorTooLargePacket
		}
	} else {
		doSend = true
	}

	if doSend {
		if e = c.sendPacketNoLock(headerData, data, cipherAdditionalData); e == nil {
			c.bytesSendReceivedEnc += len(data)
			if c.bytesSendReceivedEnc >= RekeyAfterByte {
				e = c.rekey(true, false)
			}
		}
	}
	return
}

func (c *Conn) ReadPacket(cipherAdditionalData []byte) (data []byte, e error) {
	c.r.Lock()
	defer c.r.Unlock()

	for i := 0; i < 2 && e == nil; i++ {
		var header byte
		if header, data, e = c.readPacketNoLock(cipherAdditionalData); e == nil {
			if header == headerRekey {
				if i == 0 {
					e = c.rekey(false, true)
				} else {
					e = ErrorRecursiveRekey
				}
			} else {
				break
			}
		}
	}

	return
}
