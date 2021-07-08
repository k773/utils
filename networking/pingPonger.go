package networking

import (
	"bytes"
	"time"
)

func (c *connection) StartPinging() {
	for range time.NewTicker(c.s.PingInterval).C {
		//if t.Sub(c.getLastPongReceived()) > c.s.PingTimeout {
		//	c.Close(errors.New("ping timeout"))
		//	break
		//}

		if e := c.SendPing(); e != nil {
			c.Close(e)
			break
		} else {
			c.Lock()
			c.lastPingSent = time.Now()
			c.Unlock()
		}
	}
}

func (c *connection) SendPing() (e error) {
	return c.Write(0, nil)
}

func (c *connection) SendPong() (e error) {
	return c.Write(1, nil)
}

func (c *connection) OnPing(data *bytes.Buffer) (e error) {
	return c.SendPong()
}

func (c *connection) OnPong(data *bytes.Buffer) {
	c.Lock()
	c.lastPongReceived = time.Now()
	a := c.lastPongReceived.Sub(c.lastPingSent)
	c.ping = a
	c.Unlock()
	c.cb.OnPong(a)
}
