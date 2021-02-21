package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *St) router() http.Handler {
	r := mux.NewRouter()

	mrs := a.mwfRequestSession

	r.HandleFunc("/register", mrs(a.hRegister)).Methods("GET")
	r.HandleFunc("/send", a.hSend).Methods("POST")
	r.HandleFunc("/connection_count", a.hGetConnectionCount).Methods("GET")

	return a.middleware(r)
}
