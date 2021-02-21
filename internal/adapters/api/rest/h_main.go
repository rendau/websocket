package rest

import (
	"net/http"

	"github.com/rendau/websocket/internal/domain/entities"
)

func (a *St) hRegister(w http.ResponseWriter, r *http.Request) {
	err := a.cr.ConRegister(r.Context(), w, r)
	uHandleCoreErr(err, w)
}

func (a *St) hSend(w http.ResponseWriter, r *http.Request) {
	reqObj := &entities.SendPars{}
	if !uParseRequestJSON(w, r, reqObj) {
		return
	}

	a.cr.Send(reqObj)

	w.WriteHeader(200)
}

func (a *St) hGetConnectionCount(w http.ResponseWriter, r *http.Request) {
	cnt := a.cr.GetConnectionCount()

	uRespondJSON(w, map[string]int64{
		"value": cnt,
	})
}
