package networking

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"time"
)

type Callbacks interface {
	OnMessage(msgID byte, data *bytes.Buffer)
	OnClose(error)
	OnPong(time.Duration)
}

type Settings struct {
	PingInterval time.Duration
	PingTimeout  time.Duration
}

type connection struct {
	sync.RWMutex
	r sync.Mutex
	w sync.Mutex

	c                net.Conn
	cb               Callbacks
	s                *Settings
	lastPingSent     time.Time
	lastPongReceived time.Time
	ping             time.Duration

	State string
}

type Connection interface {
	GetPing() time.Duration
	Status() string
	IsActive() bool
	Close(error)
	Run()
	Write(msgID byte, data []byte) error
}

func New(conn net.Conn, callbacks Callbacks, settings *Settings) Connection {
	if settings == nil {
		settings = &Settings{PingInterval: time.Second, PingTimeout: 5 * time.Second}
	}
	return &connection{
		c:                conn,
		cb:               callbacks,
		s:                settings,
		lastPingSent:     time.Now(),
		lastPongReceived: time.Now(),
		State:            "OPENED",
	}
}

func (c *connection) Status() string {
	c.RLock()
	defer c.RUnlock()
	return c.State
}

func (c *connection) IsActive() bool {
	return c.Status() == "OPENED"
}

func (c *connection) Close(e error) {
	c.Lock()
	defer c.Unlock()
	_ = c.c.Close()
	if c.State == "OPENED" {
		c.cb.OnClose(e)
	}
	c.State = "CLOSED"
}

func (c *connection) getLastPingSent() time.Time {
	c.RLock()
	defer c.RUnlock()
	return c.lastPingSent
}

func (c *connection) getLastPongReceived() time.Time {
	c.RLock()
	defer c.RUnlock()
	return c.lastPongReceived
}

func (c *connection) GetPing() time.Duration {
	c.RLock()
	defer c.RUnlock()
	return c.ping
}

func (c *connection) Run() {
	go c.StartPinging()

	var e error
	for e == nil {
		var msgID byte
		var buf *bytes.Buffer
		if msgID, buf, e = c.readPacket(); e == nil {
			println(msgID)
			// Ping-pong
			if msgID == 0 || msgID == 1 {
				if msgID == 0 { // Ping received
					e = c.OnPing(buf)
				} else {
					c.OnPong(buf)
				}
				continue
			}
			c.cb.OnMessage(msgID, buf)
		}
	}
	c.Close(e)
	println("Run(): finished")
}

func (c *connection) readPacket() (msgID byte, buf *bytes.Buffer, e error) {
	c.r.Lock()
	defer c.r.Unlock()

	var data []byte
	if data, e = c.read(4); e == nil {
		l := int(binary.BigEndian.Uint32(data))
		if l == 0 {
			return 0, nil, errors.New("zero packet length received - not allowed")
		}

		if data, e = c.read(l); e == nil {
			return data[0], bytes.NewBuffer(data[1:]), nil
		}
	}
	return
}

func (c *connection) read(length int) (buf []byte, e error) {
	buf = make([]byte, length)
	read := 0
	for e == nil && read != length {
		var n int
		n, e = c.c.Read(buf[read:length])
		read += n
	}
	return
}

func (c *connection) Write(msgID byte, data []byte) (e error) {
	c.w.Lock()
	defer c.w.Unlock()

	var data_ = make([]byte, 4+len(data)+1)
	binary.BigEndian.PutUint32(data_, uint32(len(data)+1))
	data_[4] = msgID
	copy(data_[5:], data)
	_, e = c.c.Write(data_)
	return
}

func (c *connection) WriteBuf(msgID byte, data *bytes.Buffer) (e error) {
	return c.Write(msgID, data.Bytes())
}
