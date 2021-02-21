package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/rendau/websocket/internal/domain/core"
	"github.com/rendau/websocket/internal/interfaces"
)

type St struct {
	lg         interfaces.Logger
	listen     string
	usrAuthUrl string
	cr         *core.St

	server *http.Server
	lChan  chan error
}

func New(
	lg interfaces.Logger,
	listen string,
	usrAuthUrl string,
	cr *core.St,
) *St {
	api := &St{
		lg:         lg,
		listen:     listen,
		usrAuthUrl: usrAuthUrl,
		cr:         cr,
		lChan:      make(chan error, 1),
	}

	api.server = &http.Server{
		Addr:         listen,
		Handler:      api.router(),
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 20 * time.Second,
	}

	return api
}

func (a *St) Start() {
	go func() {
		err := a.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			a.lg.Errorw("Http server closed", err)
			a.lChan <- err
		}
	}()
}

func (a *St) Wait() <-chan error {
	return a.lChan
}

func (a *St) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
