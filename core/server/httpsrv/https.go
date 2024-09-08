package httpsrv

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Handler func(ctx context.Context) (http.Handler, error)

type Http struct {
	ListenOn string `yaml:"listenOn"`
}

type HttpSrv struct {
	listenOn   string
	fn         Handler
	timeout    time.Duration
	httpServer *http.Server
	once       sync.Once
}
