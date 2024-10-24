package http_server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
)

func New(router http.Handler, cert *tls.Certificate, listenerName string) *listener {
	contract := &listener{
		name: listenerName,
	}

	var certificates []tls.Certificate
	if cert != nil {
		certificates = []tls.Certificate{*cert}
	}
	contract.router = &http.Server{
		TLSConfig: &tls.Config{
			Certificates: certificates,
		},
		Handler: router,
	}

	return contract
}

type listener struct {
	name   string
	router *http.Server
}

func (c *listener) Info() string {
	return fmt.Sprintf("%s, http(s) listener", c.name)
}

func (c *listener) Run(port int) (err error) {
	c.router.Addr = fmt.Sprintf(":%d", port)
	if len(c.router.TLSConfig.Certificates) > 0 {
		err = c.router.ListenAndServeTLS("", "")
	} else {
		err = c.router.ListenAndServe()
	}

	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}

	return
}

func (c *listener) Stop() error {
	return c.router.Shutdown(context.Background())
}
