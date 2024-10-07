package resolver

import (
	"context"

	"github.com/tiny-sky/Tdtm/core/registry"
	"github.com/tiny-sky/Tdtm/core/server/endpoint"
	"github.com/tiny-sky/Tdtm/log"
	"google.golang.org/grpc/resolver"
)

type defaultResolver struct {
	clientconn resolver.ClientConn
	watcher    registry.Watcher
	ctx        context.Context
	cancel     func()
}

func NewDefaultResolver(ctx context.Context, cc resolver.ClientConn, w registry.Watcher) *defaultResolver {
	r := &defaultResolver{
		clientconn: cc,
		watcher:    w,
	}
	r.ctx, r.cancel = context.WithCancel(ctx)
	return r
}

func (r *defaultResolver) watch() {
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
		}
		instances, err := r.watcher.Next()
		if err != nil {
			log.Errorf("[defaultResolver] watch err:%v", err)
			return
		}
		r.updateState(instances)
	}
}

func (r *defaultResolver) updateState(list []*registry.Instance) {
	var (
		state resolver.State
	)
	log.Infof("[defaultResolver]updateState:%v", list)

	for _, instance := range list {
		e, err := endpoint.GetHostByEndpoint(instance.Nodes, "grpc")
		if err != nil {
			log.Errorf("[updateState]GetHostByEndpoint err:%v", err)
			continue
		}
		if e == "" {
			continue
		}

		state.Addresses = append(state.Addresses, resolver.Address{
			Addr:       e,
			ServerName: instance.Name,
		})
	}
	if len(state.Addresses) == 0 {
		return
	}

	err := r.clientconn.UpdateState(state)
	if err != nil {
		log.Errorf("[updateState]UpdateState err:%v", err)
		return
	}
}

func (r *defaultResolver) ResolveNow(options resolver.ResolveNowOptions) {}

func (r *defaultResolver) Close() {
	r.cancel()
	if err := r.watcher.Stop(); err != nil {
		log.Errorf("defaultResolver close:%v", err)
	}
}
