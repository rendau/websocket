package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/rendau/websocket/internal/domain/errs"
)

func uSetContentTypeJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func uRespondJSON(w http.ResponseWriter, obj interface{}) {
	uSetContentTypeJSON(w)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		log.Panicln("Fail to encode json obj", err)
	}
}

func uParseRequestJSON(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dst)
	if err != nil {
		uRespondJSON(w, ErrRepSt{
			ErrorCode: "bad_json",
		})
		return false
	}
	return true
}

func uHandleCoreErr(err error, w http.ResponseWriter) bool {
	if err != nil {
		switch cErr := err.(type) {
		case errs.Err:
			uRespondJSON(w, ErrRepSt{
				ErrorCode: cErr.Error(),
			})
		default:
			if err != context.Canceled &&
				err != context.DeadlineExceeded {
				uRespondJSON(w, ErrRepSt{
					ErrorCode: errs.ServiceNA.Error(),
				})
			}
		}
		return true
	}

	return false
}
