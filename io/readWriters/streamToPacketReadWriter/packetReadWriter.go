package streamToPacketReadWriter

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"sync"
)

type ReadWriter struct {
	RW           io.ReadWriter
	rLock, wLock sync.Mutex

	// MaxPacketSize defines max packet size that will be handled by the application.
	// If received packet size > MaxPacketSize then an error will be returned.
	// MaxPacketSize limits only ReadPacket().
	// If provided MaxPacketSize = 0 then no limit will be applied.
	MaxPacketSize uint32
}

func NewReadWriter(rw io.ReadWriter, maxPacketSize uint32) *ReadWriter {
	return &ReadWriter{RW: rw, MaxPacketSize: maxPacketSize}
}

var ErrorPacketSizeExceedsLimit = errors.New("received packet size exceeds the limit")

// ReadPacket - dst may be nil -> a new buffer will be allocated
func (rw *ReadWriter) ReadPacket(dst *bytes.Buffer) (packet *bytes.Buffer, e error) {
	rw.rLock.Lock()
	defer rw.rLock.Unlock()

	// For some reason binary.Read is very slow
	var sizeB = make([]byte, 4)
	if _, e = io.ReadFull(rw.RW, sizeB); e == nil {
		var size = uint32(sizeB[0]) | uint32(sizeB[1])<<8 | uint32(sizeB[2])<<16 | uint32(sizeB[3])<<24
		if rw.MaxPacketSize != 0 && size > rw.MaxPacketSize {
			return nil, ErrorPacketSizeExceedsLimit
		}

		if dst == nil {
			packet = bytes.NewBuffer(make([]byte, size))
			packet.Reset()
		} else {
			packet = dst
		}
		// This is slightly faster than using io.CopyN
		var buf = make([]byte, size)
		if _, e = io.ReadFull(rw.RW, buf); e == nil {
			packet.Write(buf)
		}
	}
	return
}

func (rw *ReadWriter) WritePacket(data *bytes.Buffer) (e error) {
	rw.wLock.Lock()
	defer rw.wLock.Unlock()

	// This is a bit faster than writing with a binary.Write & calling RW.Write multiple times
	l := data.Len()
	var c = make([]byte, l+4)
	copy(c[4:], data.Bytes())
	c[0] = byte(l)
	c[1] = byte(l >> 8)
	c[2] = byte(l >> 16)
	c[3] = byte(l >> 24)

	_, e = rw.RW.Write(c)
	return
}
