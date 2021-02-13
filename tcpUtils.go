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
	Conn     net.Conn
	IsOpened bool
	Buf      struct {
		M      map[byte][][]byte
		Notify *sync.Cond
	}
	Logger log.Logger
}

var (
	ConnErrorConnectionClosed                           = errors.New("connection closed")
	ConnErrorEOF                                        = errors.New("EOF")
	ConnErrorBufferIsNotInitializedOrConnectionIsClosed = errors.New("buffer is not initialized or connection is closed")
)

func (c *ConnTools) SendMessage(msgID byte, data []byte) error {
	var bl = make([]byte, 4)
	binary.LittleEndian.PutUint32(bl, uint32(len(data)+1))
	_, err := c.Conn.Write(append(bl, append([]byte{msgID}, data...)...))
	return err
}

func (c *ConnTools) StartMonitoring() error {
	c.Buf.M = make(map[byte][][]byte)
	c.Buf.Notify = sync.NewCond(&sync.Mutex{})
	defer c.Buf.Notify.Broadcast()

	for range time.NewTicker(10 * time.Millisecond).C {
		msgID, data, err := c.ReadMessageFromConnection()
		if !c.IsOpened {
			return ConnErrorConnectionClosed
		}
		if err != nil && err.Error() == "EOF" {
			c.Close("error while reading: EOF")
			return ConnErrorEOF
		} else if len(data) != 0 {
			c.Buf.Notify.L.Lock()
			c.Buf.M[msgID] = append(c.Buf.M[msgID], data)
			c.Buf.Notify.Broadcast()
			c.Buf.Notify.L.Unlock()
		}
	}
	return nil
}

func (c *ConnTools) ReadMessageFromBuffer(msgID byte) ([]byte, error) {
	var data []byte

	c.Buf.Notify.L.Lock()
	defer c.Buf.Notify.L.Unlock()
	for {
		if c.Buf.M == nil || !c.IsOpened {
			return nil, ConnErrorBufferIsNotInitializedOrConnectionIsClosed
		}

		if c.Buf.M[msgID] != nil && len(c.Buf.M[msgID]) != 0 {
			data = c.Buf.M[msgID][0]
			c.Buf.M[msgID] = c.Buf.M[msgID][1:]
		}

		if len(data) != 0 {
			break
		}
		c.Buf.Notify.Wait()
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

	if _, err = c.Conn.Read(msgLengthB); err == nil {
		msgLength := int64(binary.LittleEndian.Uint32(msgLengthB))

		var buf = make([]byte, 1024)
	receive:
		var r int
		r, err = c.Conn.Read(buf)
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
	_ = c.Conn.Close()
	c.IsOpened = false
	c.Logger.Println(c.Conn.RemoteAddr().String(), "closed:", reason)
}
