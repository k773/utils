package aesConn

import (
	"crypto/aes"
	"github.com/k773/utils"
	"github.com/k773/utils/io/conn/encryptedConn"
	"net"
	"time"
)

func Dial(addr string, aesKey []byte, hskTimeout time.Duration) (_ net.Conn, e error) {
	conn, e := net.Dial("tcp", addr)
	if e != nil {
		return
	}
	if hskTimeout != 0 {
		_ = conn.SetDeadline(time.Now().Add(hskTimeout))
		defer conn.SetDeadline(time.Time{})
	}
	return UpgradeClient(conn, aesKey)
}

func DialTimeout(addr string, aesKey []byte, dialTimeout, hskTimeout time.Duration) (_ net.Conn, e error) {
	conn, e := net.DialTimeout("tcp", addr, dialTimeout)
	if e != nil {
		return
	}
	if hskTimeout != 0 {
		_ = conn.SetDeadline(time.Now().Add(hskTimeout))
		defer conn.SetDeadline(time.Time{})
	}
	return UpgradeClient(conn, aesKey)
}

func UpgradeClient(conn net.Conn, aesKey []byte) (_ *encryptedConn.Conn, e error) {
	var c = &encryptedConn.Conn{Conn: conn}

	// Client and server already share the same encryption key
	// Handshake:
	// server -> client: [plain/packet]: iv + 32 random bytes
	// client -> server: [enc/packet]: reverse 32 bytes
	// server -> client: [enc/packet]: <empty> / conn close

	var r []byte
	if r, e = c.ReadPacket(); e == nil {
		if e = GenerateCiphers(r[:aes.BlockSize], aesKey, c); e == nil {
			if e = c.WritePacket(utils.Reverse(r[aes.BlockSize:])); e == nil {
				_, e = c.ReadPacket()
			}
		}
	}

	return c, e
}
