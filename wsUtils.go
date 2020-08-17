package utils

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var Upgrader = websocket.Upgrader{}

type WsStruct struct {
	S        sync.Mutex
	Conn     *websocket.Conn
	M        map[[16]byte][][]byte //msgID:[]msg; map may be nil (do not use readMsgIDMessage() and readMessagesThread() then)
	IsOpened bool                  //Needs to be set manually
}

//DEPRECATED
func ReadRawMessage(conn *websocket.Conn) ([]byte, error) {
	_, response, err := conn.ReadMessage()
	if err != nil {
		_ = conn.Close()
	}
	return response, err
}

//DEPRECATED
func WriteRawMessage(conn *websocket.Conn, data []byte) error {
	err := conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		_ = conn.Close()
		return err
	}
	return nil
}

func (ws *WsStruct) ReadWsMessage() ([16]byte, []byte, error) { //Returns msgId, msg, err
	_, response, err := ws.Conn.ReadMessage()
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

func (ws *WsStruct) ReadMsgIDMessage(msgID [16]byte) ([]byte, error) {
	//fmt.Println(msgID)
	for {
		ws.S.Lock()
		if value, has := ws.M[msgID]; has && len(ws.M[msgID]) > 0 {
			//Removing received data from the map

			ws.M[msgID] = ws.M[msgID][1:]
			ws.S.Unlock()
			return value[0], nil
		}
		ws.S.Unlock()
		if !ws.IsOpened {
			return nil, errors.New("Ws closed!")
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (ws *WsStruct) WriteMessage(msgID [16]byte, msg []byte) error {
	return ws.Conn.WriteMessage(websocket.BinaryMessage, append(msgID[:], msg...))
}

func (ws *WsStruct) ReadMessagesThread() error {
	for {
		msgID, msg, err := ws.ReadWsMessage()
		if err != nil {
			return err
		}

		ws.S.Lock()
		ws.M[msgID] = append((ws.M)[msgID], msg)
		ws.S.Unlock()
	}
}
