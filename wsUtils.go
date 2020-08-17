package utils

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{}

type wsStruct struct {
	s        sync.Mutex
	conn     *websocket.Conn
	m        map[[16]byte][][]byte //msgID:[]msg; map may be nil (do not use readMsgIDMessage() and readMessagesThread() then)
	isOpened bool                  //Needs to be set manually
}

//DEPRECATED
func readRawMessage(conn *websocket.Conn) ([]byte, error) {
	_, response, err := conn.ReadMessage()
	if err != nil {
		_ = conn.Close()
	}
	return response, err
}

//DEPRECATED
func writeRawMessage(conn *websocket.Conn, data []byte) error {
	err := conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		_ = conn.Close()
		return err
	}
	return nil
}

func (ws *wsStruct) readWsMessage() ([16]byte, []byte, error) { //Returns msgId, msg, err
	_, response, err := ws.conn.ReadMessage()
	if err != nil {
		return [16]byte{}, nil, err
	}
	if len(response) < 16 {
		return [16]byte{}, nil, errors.New("Received message's length is less then msgID min length (16)")
	}

	//Cant just cast type []byte to type [16]byte
	var msgID [16]byte
	copy(msgID[:], response[:16])

	return msgID, response[16:], nil
}

func (ws *wsStruct) readMsgIDMessage(msgID [16]byte) ([]byte, error) {
	//fmt.Println(msgID)
	for {
		ws.s.Lock()
		if value, has := ws.m[msgID]; has && len(ws.m[msgID]) > 0 {
			//Removing received data from the map

			ws.m[msgID] = ws.m[msgID][1:]
			ws.s.Unlock()
			return value[0], nil
		}
		ws.s.Unlock()
		if !ws.isOpened {
			return nil, errors.New("Ws closed!")
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (ws *wsStruct) writeMessage(msgID [16]byte, msg []byte) error {
	return ws.conn.WriteMessage(websocket.BinaryMessage, append(msgID[:], msg...))
}

func (ws *wsStruct) readMessagesThread() error {
	for {
		msgID, msg, err := ws.readWsMessage()
		if err != nil {
			return err
		}

		ws.s.Lock()
		ws.m[msgID] = append((ws.m)[msgID], msg)
		ws.s.Unlock()
	}
}
