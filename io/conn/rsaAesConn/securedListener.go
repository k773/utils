package rsaAesConn

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"github.com/k773/utils/io/conn/aesConn"
	"github.com/k773/utils/io/conn/encryptedConn"
	"net"
	"time"
)

type Listener struct {
	net.Listener
	HandshakeTimeout time.Duration

	AesKey  []byte
	Private *rsa.PrivateKey

	// BeforeEncInit is called after an underlying listtener has accepted new connection.
	// If returned error != nil, connection will not be processed but immediately closed.
	BeforeEncInit func(conn net.Conn) error
	// OnEncInitError is called if there was an error during ecryption initialization phase.
	// Error generated by BeforeEncInit will not be passed th this func.
	OnEncInitError func(conn net.Conn, e error)
}

// Accept accepts new connection and encrypts it before returning.
// Any error except underlying listener's won't be returned, instead use Listener.OnEncInitError to handle it.
func (l *Listener) Accept() (net.Conn, error) {
	var conn net.Conn
	var e error
	for e == nil {
		conn, e = l.Listener.Accept()
		if e == nil {
			if e = l.BeforeEncInit(conn); e == nil {
				// Limitting handshake time
				if l.HandshakeTimeout != 0 {
					_ = conn.SetDeadline(time.Now().Add(l.HandshakeTimeout))
				}

				// Handshaking
				if conn, e = UpgradeServer(conn, l.AesKey, l.Private); e == nil {
					conn.SetDeadline(time.Time{})
					break
				}

				// Reporting connection error
				if e != nil {
					l.OnEncInitError(conn, e)
				}
			}
			// Resetting error; any error during encryption initialization phase should be handled by l.OnEncInitError
			if e != nil {
				_ = conn.Close()
				e = nil
			}
		}
	}
	return conn, e
}

func UpgradeServer(conn net.Conn, aesKey []byte, private *rsa.PrivateKey) (_ net.Conn, e error) {
	//return conn, t, nil
	var c = &encryptedConn.Conn{Conn: conn}

	// Client and server already share the same encryption key
	// Handshake:
	// client -> server: [rsa/packet]: iv + temp aes key
	// server -> client: [temp_aes/packet]: iv + aes key + 32 random bytes
	// Handshake by aesConn.UpgradeServer:
	// client -> server: [end_aes/packet]: reversed 32 bytes
	// server -> client: [end_aes/packet]: <empty> / conn close

	var initKey []byte
	if initKey, e = c.ReadPacket(); e == nil {
		if initKey, e = rsa.DecryptOAEP(sha512.New(), rand.Reader, private, initKey, []byte("handshake")); e == nil {
			if e = aesConn.GenerateCiphers(initKey, nil, c); e == nil {
				c, e = aesConn.UpgradeServer(c, aesKey)
			}
		}
	}

	return c, e
}
