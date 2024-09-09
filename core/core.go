package core

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/tiny-sky/Tdtm/core/endpoint"
	"github.com/tiny-sky/Tdtm/core/registry"
	"github.com/tiny-sky/Tdtm/log"
	"golang.org/x/sync/errgroup"
)

type Option func(core *Core)

type Core struct {
	server       []Server
	stopCtx      context.Context
	runWaitGroup sync.WaitGroup
	errGroup     *errgroup.Group
	cancel       func()
	once         sync.Once
	registry     registry.Registry
	instance     *registry.Instance
}

func WithServers(srvs ...Server) Option {
	return func(core *Core) {
		core.server = append(core.server, srvs...)
	}
}

func WithRegistry(r registry.Registry) Option {
	return func(core *Core) {
		core.registry = r
	}
}

func New(opts ...Option) *Core {
	core := &Core{
		runWaitGroup: sync.WaitGroup{},
		once:         sync.Once{},
	}
	for _, opt := range opts {
		opt(core)
	}
	return core
}

func (core *Core) Run(ctx context.Context) error {
	var c1 context.Context

	c1, core.cancel = context.WithCancel(ctx)
	core.errGroup, core.stopCtx = errgroup.WithContext(c1)

	core.instance = registry.NewInstance()

	for _, server := range core.server {
		core.runWaitGroup.Add(1)

		if e, ok := server.(endpoint.Endpoint); ok {
			core.instance.Nodes = append(core.instance.Nodes, e.Endpoint().String())
		}
		srv := server
		core.errGroup.Go(func() error {
			<-core.stopCtx.Done()
			return srv.Stop(ctx)
		})
	}

	core.runWaitGroup.Wait()
	log.Infof("start")

	if core.registry != nil {
		if err := core.registry.Register(c1, core.instance); err != nil {
			return err
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	core.errGroup.Go(func() error {
		select {
		case <-core.stopCtx.Done():
			return core.stopCtx.Err()
		case <-c:
			return core.Stop()
		}
	})
	if err := core.errGroup.Wait(); err != nil {
		return err
	}
	return nil
}

func (core *Core) Stop() (err error) {
	if core.cancel == nil {
		return nil
	}
	core.once.Do(func() {
		if core.registry != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = core.registry.DeRegister(ctx, core.instance)
		}
		core.cancel()
	})
	return
}