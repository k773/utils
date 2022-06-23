package packetReadWriter

import (
	"bytes"
	"encoding/binary"
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

	var size uint32
	if e = binary.Read(rw.RW, binary.LittleEndian, &size); e == nil {
		if rw.MaxPacketSize != 0 && size > rw.MaxPacketSize {
			return nil, ErrorPacketSizeExceedsLimit
		}

		if dst != nil {
			packet = dst
		} else {
			packet = bytes.NewBuffer(make([]byte, size))
			packet.Reset()
		}

		_, e = io.CopyN(packet, rw.RW, int64(size))
	}
	return
}

func (rw *ReadWriter) WritePacket(data *bytes.Buffer) (e error) {
	rw.wLock.Lock()
	defer rw.wLock.Unlock()

	if e = binary.Write(rw.RW, binary.LittleEndian, uint32(data.Len())); e == nil {
		//_, e = io.Copy(rw.RW, data)
		_, e = rw.RW.Write(data.Bytes())
	}
	return
}
