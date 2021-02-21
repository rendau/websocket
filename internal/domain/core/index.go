package core

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rendau/websocket/internal/domain/entities"
	"github.com/rendau/websocket/internal/domain/errs"
	"github.com/rendau/websocket/internal/interfaces"
)

const sessionContextKey = "usr_session"

type St struct {
	lg interfaces.Logger

	wsUpgrader websocket.Upgrader
	cons       []*entities.ConSt
	consMU     sync.RWMutex
}

func New(lg interfaces.Logger) (*St, error) {
	c := &St{
		lg: lg,
		wsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		cons:   make([]*entities.ConSt, 0),
		consMU: sync.RWMutex{},
	}

	return c, nil
}

func (c *St) ContextWithSession(ctx context.Context, ses *entities.Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, ses)
}

func (c *St) ContextGetSession(ctx context.Context) *entities.Session {
	contextV := ctx.Value(sessionContextKey)
	if contextV == nil {
		return &entities.Session{}
	}

	switch ses := contextV.(type) {
	case *entities.Session:
		return ses
	default:
		c.lg.Fatal("wrong type of session in context")
		return nil
	}
}

func (c *St) SesRequireAuth(ses *entities.Session) error {
	if ses.ID == 0 {
		return errs.NotAuthorized
	}
	return nil
}
