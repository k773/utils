package rsaAesConn

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"github.com/k773/utils/io/conn/aesConn"
	"github.com/k773/utils/io/conn/encryptedConn"
	"net"
	"time"
)

func Dial(addr string, aesKey []byte, serverPubKey *rsa.PublicKey, hskTimeout time.Duration) (_ net.Conn, e error) {
	conn, e := net.Dial("tcp", addr)
	if e != nil {
		return
	}
	if hskTimeout != 0 {
		_ = conn.SetDeadline(time.Now().Add(hskTimeout))
		defer conn.SetDeadline(time.Time{})
	}
	return UpgradeClient(conn, aesKey, serverPubKey)
}

func DialTimeout(addr string, aesKey []byte, serverPubKey *rsa.PublicKey, dialTimeout, hskTimeout time.Duration) (_ net.Conn, e error) {
	conn, e := net.DialTimeout("tcp", addr, dialTimeout)
	if e != nil {
		return
	}
	if hskTimeout != 0 {
		_ = conn.SetDeadline(time.Now().Add(hskTimeout))
		defer conn.SetDeadline(time.Time{})
	}
	return UpgradeClient(conn, aesKey, serverPubKey)
}

func UpgradeClient(conn net.Conn, aesKey []byte, serverPubKey *rsa.PublicKey) (_ *encryptedConn.Conn, e error) {
	var c = &encryptedConn.Conn{Conn: conn}

	// Client and server already share the same encryption key
	// Handshake:
	// client -> server: [rsa/packet]: iv + temp aes key
	// Handshake by aesConn.UpgradeClient:
	// server -> client: [temp_aes/packet]: iv + aes key + 32 random bytes
	// client -> server: [end_aes/packet]: reversed 32 bytes
	// server -> client: [end_aes/packet]: <empty> / conn close

	var initKey = make([]byte, aes.BlockSize+aes.BlockSize)
	var r []byte
	if _, e = rand.Read(initKey); e == nil {
		if r, e = rsa.EncryptOAEP(sha512.New(), rand.Reader, serverPubKey, initKey, []byte("handshake")); e == nil {
			if e = c.WritePacket(r); e == nil {
				if e = aesConn.GenerateCiphers(initKey, nil, c); e == nil {
					return aesConn.UpgradeClient(c, aesKey)
				}
			}
		}
	}

	return c, e
}
