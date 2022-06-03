package packetedConn

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/k773/utils"
	"io"
	"net"
	"sync"
)

var buffPool = &sync.Pool{
	New: func() any { return bytes.NewBuffer(nil) },
}

// Packet structure: [4 bytes: message size encoded in BE excluding this header; message]

// Conn implements io.ReadWriter, however it is recommended to use utils.PacketReadWriter instead
type Conn struct {
	net.Conn
	readBuf *bytes.Buffer

	ReadPacketSizeLimit uint32

	r, w sync.Mutex
}

// Write is just a wrapper around ReadPacket implementing io.Writer. Use Conn.WritePacket when possible
func (c *Conn) Write(src []byte) (n int, e error) {
	// No need in using buffers pool as src is already allocated
	e = c.WritePacket(bytes.NewBuffer(src))
	return utils.If(e == nil, len(src), 0), e
}

// Read is just a wrapper around ReadPacket implementing io.Reader. Use Conn.ReadPacket when possible
func (c *Conn) Read(dst []byte) (n int, e error) {
	if c.readBuf.Len() == 0 {
		n, e = c.readBuf.Read(dst)
		return
	}
	var packet *bytes.Buffer
	if packet, e = c.ReadPacket(); e == nil {
		if _, e = c.readBuf.ReadFrom(packet); e == nil {
			n, e = c.readBuf.Read(dst)
		}
	}
	ReleaseBuffer(packet)
	return
}

func (c *Conn) Close() error {
	return c.Conn.Close()
}

// ReadPacket returns a buffer that must be returned to the pool by a caller
func (c *Conn) ReadPacket() (*bytes.Buffer, error) {
	c.w.Lock()
	defer c.w.Unlock()

	return ReadPacket(c.Conn, c.ReadPacketSizeLimit)
}

// WritePacket returns a buffer that must be returned to the pool by a caller
func (c *Conn) WritePacket(src *bytes.Buffer) error {
	c.w.Lock()
	defer c.w.Unlock()

	return WritePacket(c.Conn, src)
}

func (c *Conn) GetBuffer(size int) *bytes.Buffer {
	return GetBuffer(size)
}

func (c *Conn) ReleaseBuffer(buf *bytes.Buffer) {
	ReleaseBuffer(buf)
}

// WritePacket will not return the given buffer to the pool
func WritePacket(c net.Conn, src *bytes.Buffer) error {
	e := binary.Write(c, binary.LittleEndian, uint32(src.Len()))
	if e == nil {
		_, e = c.Write(src.Bytes())
	}
	return e
}

// MustRead returns a buffer that must be returned to the pool by a caller
func MustRead(c net.Conn, n int) (buf *bytes.Buffer, e error) {
	buf = GetBuffer(n)
	_, e = io.CopyN(buf, c, int64(n))
	return
}

// ReadPacket returns a buffer that must be returned to the pool by a caller
func ReadPacket(c net.Conn, packetSizeLimit uint32) (buf *bytes.Buffer, e error) {
	sizeBuf, e := MustRead(c, 4)
	defer ReleaseBuffer(sizeBuf)
	if e == nil {
		var size = int(binary.LittleEndian.Uint32(sizeBuf.Bytes()))
		if packetSizeLimit > 0 && size > int(packetSizeLimit) {
			_ = c.Close()
			return nil, errors.New("read packet size limite exceeded")
		} else {
			buf, e = MustRead(c, size)
		}
	}
	return
}

/*
	Buffering is disabled due to very large memory consumption
*/

func GetBuffer(size int) *bytes.Buffer {
	b := bytes.NewBuffer(nil)
	b.Grow(size)
	return b
	var buf = buffPool.Get().(*bytes.Buffer)
	buf.Grow(size)
	return buf
}

func ReleaseBuffer(buf *bytes.Buffer) {

	return
	if buf != nil {
		buf.Reset()
		buffPool.Put(buf)
	}
}

// CopyN is an exact copy of io.CopyN except io.Copy is replaced with io.CopyBuffer
func CopyN(dst *bytes.Buffer, src io.Reader, n int64) (written int64, err error) {
	//var tmp = GetBuffer(int(n))
	dst.Grow(int(n))
	written, err = io.CopyBuffer(dst, io.LimitReader(src, n), dst.Bytes()[:n]) // [:n] - expanding slice to needed capacity
	//ReleaseBuffer(tmp)
	if written == n {
		return n, nil
	}
	if written < n && err == nil {
		// src stopped early; must have been EOF.
		err = io.EOF
	}
	return
}
