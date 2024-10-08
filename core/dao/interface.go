package dao

import (
	"context"
	"github.com/tiny-sky/Tdtm/core/consts"
	"github.com/tiny-sky/Tdtm/core/dao/entity"
)

type TransactionDao interface {
	BranchDao
	GlobalDao
}

type BranchDao interface {
	CreateBatches(ctx context.Context, list entity.BranchList) error
	GetBranches(ctx context.Context, gid string) (entity.BranchList, error)
	UpdateBranchStateByGid(ctx context.Context, branchId string,
		state consts.BranchState, errMsg string) (int64, error)
}

type GlobalDao interface {
	FindProcessingList(ctx context.Context, limit, maxTimes int) (list []*entity.Global, err error)
	IncrTryTimes(ctx context.Context, gid string, nextCronTime int) error
	CreateGlobal(ctx context.Context, global *entity.Global) error
	GetGlobal(ctx context.Context, gid string) (entity.Global, error)
	UpdateGlobalStateByGid(ctx context.Context, gid string,
		state consts.GlobalState) (int64, error)
}

var dao TransactionDao

func GetTransaction() TransactionDao {
	return dao
}

func SetTransaction(d TransactionDao) {
	if d == nil {
		panic("TransactionDao must no be nil")
	}
	dao = d
}
