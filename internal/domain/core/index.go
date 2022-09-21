package core

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rendau/dop/adapters/jwk"
	"github.com/rendau/dop/adapters/logger"
	"github.com/rendau/websocket/internal/domain/types"
)

const sessionContextKey = "usr_session"

type St struct {
	lg  logger.Lite
	jwk jwk.Jwk

	wsUpgrader websocket.Upgrader

	cons   map[int64]*types.ConSt
	consMU sync.RWMutex
}

func New(lg logger.Lite, jwk jwk.Jwk) *St {
	return &St{
		lg:  lg,
		jwk: jwk,

		wsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},

		cons:   make(map[int64]*types.ConSt, 50),
		consMU: sync.RWMutex{},
	}
}
