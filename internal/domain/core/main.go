package core

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rendau/dop/adapters/jwt"
	"github.com/rendau/dop/dopErrs"
	"github.com/rendau/websocket/internal/domain/types"
)

func (c *St) ConRegister(w http.ResponseWriter, r *http.Request, token string) error {
	var err error

	tokenValid, err := c.jwk.Validate(token)
	if err != nil {
		return err
	}
	if !tokenValid {
		return dopErrs.NotAuthorized
	}

	jwtPayload := &types.JwtPayload{}

	err = jwt.ParsePayload(token, jwtPayload)
	if err != nil {
		return dopErrs.NotAuthorized
	}

	usrId, _ := strconv.ParseInt(jwtPayload.Sub, 10, 64)
	if usrId <= 0 {
		return dopErrs.NotAuthorized
	}

	wsCon, err := c.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		c.lg.Errorw("Fail to upgrade websocket", err)
		return dopErrs.ServiceNA
	}

	c.consMU.Lock()
	defer c.consMU.Unlock()

	con := &types.ConSt{
		Con:   wsCon,
		UsrId: usrId,
		Out:   make(chan []byte, 50),
	}

	c.cons[usrId] = con

	go c.conWriter(con)
	go c.conReader(con)

	return nil
}

func (c *St) Send(pars *types.SendReqSt) {
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

func (c *St) GetConnectionCount() types.ConnectionCountRepSt {
	c.consMU.RLock()
	defer c.consMU.RUnlock()

	return types.ConnectionCountRepSt{Value: int64(len(c.cons))}
}

func (c *St) conUnregister(con *types.ConSt) {
	c.consMU.Lock()
	defer c.consMU.Unlock()

	_ = con.Con.Close()

	usrId := con.UsrId

	if con := c.cons[usrId]; con != nil {
		close(con.Out)
		delete(c.cons, usrId)
	}
}

func (c *St) conWriter(con *types.ConSt) {
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

func (c *St) conWriteMsg(con *types.ConSt, tp int, msg []byte) error {
	err := con.Con.SetWriteDeadline(time.Now().Add(30 * time.Second))
	if err != nil {
		c.lg.Errorw("Fail to set timeout for ws-write", err)
		return err
	}

	return con.Con.WriteMessage(tp, msg)
}

func (c *St) conReader(con *types.ConSt) {
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
		// c.lg.Infow("Read")
	}
}
