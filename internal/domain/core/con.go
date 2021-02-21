package core

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rendau/websocket/internal/domain/entities"
	"github.com/rendau/websocket/internal/domain/errs"
)

func (c *St) ConRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var err error

	ses := c.ContextGetSession(ctx)
	if err = c.SesRequireAuth(ses); err != nil {
		return err
	}

	wsCon, err := c.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		c.lg.Errorw("Fail to upgrade websocket", err)
		return errs.ServiceNA
	}

	c.consMU.Lock()
	defer c.consMU.Unlock()

	con := &entities.ConSt{
		Con:   wsCon,
		UsrId: ses.ID,
		Out:   make(chan []byte, 50),
	}

	c.cons = append(c.cons, con)

	go c.conWriter(con)
	go c.conReader(con)

	return nil
}

func (c *St) Send(pars *entities.SendPars) {
	if len(pars.UsrIds) == 0 {
		return
	}

	c.consMU.RLock()
	defer c.consMU.RUnlock()

	for _, con := range c.cons {
		for _, usrId := range pars.UsrIds {
			if con.UsrId == usrId {
				con.Out <- pars.Message
			}
		}
	}
}

func (c *St) GetConnectionCount() int64 {
	c.consMU.RLock()
	defer c.consMU.RUnlock()

	return int64(len(c.cons))
}

func (c *St) conUnregister(con *entities.ConSt) {
	c.consMU.Lock()
	defer c.consMU.Unlock()

	_ = con.Con.Close()

	newCons := make([]*entities.ConSt, 0, len(c.cons))

	for _, c := range c.cons {
		if c == con {
			close(c.Out)
		} else {
			newCons = append(newCons, c)
		}
	}

	c.cons = newCons
}

func (c *St) conWriter(con *entities.ConSt) {
	var err error
	var msg []byte
	var ok bool

	ticker := time.NewTicker(20 * time.Second)

	defer func() {
		ticker.Stop()
		c.conUnregister(con)
	}()

	for {
		select {
		case msg, ok = <-con.Out:
			if !ok {
				_ = c.conWriteMsg(con, websocket.CloseMessage, []byte{})
				return
			}
			if err = c.conWriteMsg(con, websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err = c.conWriteMsg(con, websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *St) conWriteMsg(con *entities.ConSt, tp int, msg []byte) error {
	err := con.Con.SetWriteDeadline(time.Now().Add(30 * time.Second))
	if err != nil {
		c.lg.Errorw("Fail to set timeout for ws-write", err)
		return err
	}

	return con.Con.WriteMessage(tp, msg)
}

func (c *St) conReader(con *entities.ConSt) {
	var err error

	defer c.conUnregister(con)

	con.Con.SetCloseHandler(func(code int, text string) error {
		_ = con.Con.Close()
		return nil
	})

	for {
		if _, _, err = con.Con.ReadMessage(); err != nil {
			break
		}
		c.lg.Info("reading")
	}
}
