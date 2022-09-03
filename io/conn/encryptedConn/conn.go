package encryptedConn

import (
	"encoding/binary"
	"net"
)

type Conn struct {
	net.Conn

	Encrypt func(dst, src []byte)
	Decrypt func(dst, src []byte)
}

func (c *Conn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	if c.Decrypt != nil {
		c.Decrypt(b[:n], b[:n])
	}
	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	if c.Encrypt != nil {
		c.Encrypt(b, b)
	}
	return c.Conn.Write(b)
}

func (c *Conn) mustRead(n int) (buf []byte, e error) {
	var read = 0
	buf = make([]byte, n)
	for read < n && e == nil {
		var a int
		if a, e = c.Read(buf[read:]); e == nil {
			read += a
		}
	}
	return
}

func (c *Conn) ReadPacket() (buf []byte, e error) {
	if buf, e = c.mustRead(2); e == nil {
		buf, e = c.mustRead(int(binary.LittleEndian.Uint16(buf)))
	}
	return
}

func (c *Conn) WritePacket(raw []byte) (e error) {
	var data = make([]byte, len(raw)+2)
	binary.LittleEndian.PutUint16(data[:2], uint16(len(raw)))
	copy(data[2:], raw)

	_, e = c.Write(data)
	return
}
