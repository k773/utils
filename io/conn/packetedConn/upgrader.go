package packetedConn

import (
	"bytes"
	"net"
)

func Upgrade(c net.Conn, readPacketSizeLimit uint32) *Conn {
	return &Conn{
		Conn:                c,
		readBuf:             bytes.NewBuffer(nil),
		ReadPacketSizeLimit: readPacketSizeLimit,
	}
}
