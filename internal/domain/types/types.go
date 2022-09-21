package types

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type ConSt struct {
	Con   *websocket.Conn
	UsrId int64
	Out   chan []byte
}

type SendReqSt struct {
	UsrIds  []int64         `json:"usr_ids"`
	Message json.RawMessage `json:"message"`
}

type ConnectionCountRepSt struct {
	Value int64 `json:"value"`
}

type JwtPayload struct {
	Sub string `json:"sub"`
}
