package coordinator

import (
	"context"

	"github.com/tiny-sky/Tdtm/core/dao"
	"github.com/tiny-sky/Tdtm/core/entity"
)

type Executor interface {
	Phase1(ctx context.Context, global *entity.Global) error
	Phase2(ctx context.Context, global *entity.Global) error
	Close(ctx context.Context) error
}

type Coordinator struct {
	dao                 dao.TransactionDao
	automaticExecution2 bool
	executor            Executor
	closeFn             func(ctx context.Context) error
}

func NewCoordinator(dao dao.TransactionDao, execute Executor, automaticExecution2 bool) *Coordinator {
	c := &Coordinator{
		dao:                 dao,
		automaticExecution2: automaticExecution2,
		executor:            execute,
		closeFn:             execute.Close,
	}
	return c
}
