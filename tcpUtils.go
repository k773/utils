package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

type ConnTools struct {
	conn     net.Conn
	isOpened bool
	buf      struct {
		m      map[byte][][]byte
		notify *sync.Cond
	}
	logger log.Logger
}

var (
	ConnErrorConnectionClosed                           = errors.New("connection closed")
	ConnErrorEOF                                        = errors.New("EOF")
	ConnErrorBufferIsNotInitializedOrConnectionIsClosed = errors.New("buffer is not initialized or connection is closed")
)

func (c *ConnTools) SendMessage(msgID byte, data []byte) error {
	var bl = make([]byte, 4)
	binary.LittleEndian.PutUint32(bl, uint32(len(data)+1))
	_, err := c.conn.Write(append(bl, append([]byte{msgID}, data...)...))
	return err
}

func (c *ConnTools) StartMonitoring() error {
	c.buf.m = make(map[byte][][]byte)
	c.buf.notify = sync.NewCond(&sync.Mutex{})
	defer c.buf.notify.Broadcast()

	for range time.NewTicker(10 * time.Millisecond).C {
		msgID, data, err := c.ReadMessageFromConnection()
		if !c.isOpened {
			return ConnErrorConnectionClosed
		}
		if err != nil && err.Error() == "EOF" {
			c.Close("error while reading: EOF")
			return ConnErrorEOF
		} else if len(data) != 0 {
			c.buf.notify.L.Lock()
			c.buf.m[msgID] = append(c.buf.m[msgID], data)
			c.buf.notify.Broadcast()
			c.buf.notify.L.Unlock()
		}
	}
	return nil
}

func (c *ConnTools) ReadMessageFromBuffer(msgID byte) ([]byte, error) {
	var data []byte

	c.buf.notify.L.Lock()
	defer c.buf.notify.L.Unlock()
	for {
		if c.buf.m == nil || !c.isOpened {
			return nil, ConnErrorBufferIsNotInitializedOrConnectionIsClosed
		}

		if c.buf.m[msgID] != nil && len(c.buf.m[msgID]) != 0 {
			data = c.buf.m[msgID][0]
			c.buf.m[msgID] = c.buf.m[msgID][1:]
		}

		if len(data) != 0 {
			break
		}
		c.buf.notify.Wait()
	}

	return data, nil
}

// Msg structure: length (4bytes) + id (1byte) + data
func (c *ConnTools) ReadMessageFromConnection() (byte, []byte, error) {
	msgLengthB := make([]byte, 4)
	var err error
	var read int64
	var msgID byte
	var res bytes.Buffer

	if _, err = c.conn.Read(msgLengthB); err == nil {
		msgLength := int64(binary.LittleEndian.Uint32(msgLengthB))

		var buf = make([]byte, 1024)
	receive:
		var r int
		r, err = c.conn.Read(buf)
		read += int64(r)
		res.Write(buf[:r])
		if read < msgLength && err == nil {
			goto receive
		}
	}

	if err == nil {
		msgID, err = res.ReadByte()
	}

	return msgID, res.Bytes(), err
}

func (c *ConnTools) Close(reason string) {
	_ = c.conn.Close()
	c.isOpened = false
	c.logger.Println(c.conn.RemoteAddr().String(), "closed:", reason)
}
