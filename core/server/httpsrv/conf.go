package httpsrv

import (
	"context"
	"net/http"
)

type Http struct {
	ListenOn string `yaml:"listenOn"`
}

type Handler func(ctx context.Context) (http.Handler, error)
