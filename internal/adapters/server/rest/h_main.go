package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dopHttps "github.com/rendau/dop/adapters/server/https"
	"github.com/rendau/websocket/internal/domain/types"
)

// @Router  /register [get]
// @Tags    main
// @Success 200
// @Failure 400 {object} dopTypes.ErrRep
func (a *St) hRegister(c *gin.Context) {
	token := dopHttps.GetAuthToken(c)

	dopHttps.Error(c, a.core.ConRegister(c.Writer, c.Request, token))
}

// @Router  /send [post]
// @Tags    main
// @Param   body body types.SendReqSt false "body"
// @Success 200
func (a *St) hSend(c *gin.Context) {
	reqObj := &types.SendReqSt{}
	if !dopHttps.BindJSON(c, reqObj) {
		return
	}

	a.core.Send(reqObj)
}

// @Router  /connection_count [get]
// @Tags    main
// @Produce json
// @Success 200 {object} types.ConnectionCountRepSt
func (a *St) hGetConnectionCount(c *gin.Context) {
	c.JSON(http.StatusOK, a.core.GetConnectionCount())
}
