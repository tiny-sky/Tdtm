package transport

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/tiny-sky/Tdtm/core/transport/common"
	"github.com/tiny-sky/Tdtm/core/transport/grpc"
	http_ "github.com/tiny-sky/Tdtm/core/transport/http"
	"github.com/tiny-sky/Tdtm/log"
)

type (
	Transport interface {
		GetType() common.Net
		Request(ctx context.Context, url string, req *common.Req) (*common.Resp, error)
		Close(ctx context.Context) error
	}
	Manager interface {
		GetTransporter(net common.Net) (Transport, error)
		Close(ctx context.Context) error
	}
)

type manager struct {
	m sync.Map
}

func NewManager() *manager {
	manager := &manager{
		m: sync.Map{},
	}

	var list []Transport

	list = append(list, http_.NewTransporter(), grpc.NewTransporter())
	for _, transport := range list {
		manager.m.Store(string(transport.GetType()), transport)
	}
	return manager
}

func (manager *manager) GetTransporter(net common.Net) (Transport, error) {
	val, ok := manager.m.Load(string(net))
	if !ok {
		return nil, errors.New("not found transport")
	}
	return val.(Transport), nil
}

func (manager *manager) Close(ctx context.Context) error {
	manager.m.Range(func(key, value any) bool {
		if err := value.(Transport).Close(ctx); err != nil {
			log.Infof(fmt.Sprintf("[Manager] stop err:%v", err), "net", key, "transporter", value)
		}
		return true
	})
	return nil
}
