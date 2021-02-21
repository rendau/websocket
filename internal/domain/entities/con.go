package entities

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type ConSt struct {
	Con   *websocket.Conn
	UsrId int64
	Out   chan []byte
}

type SendPars struct {
	UsrIds  []int64         `json:"usr_ids"`
	Message json.RawMessage `json:"message"`
}
