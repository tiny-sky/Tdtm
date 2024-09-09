package mysqlx

import (
	"context"

	"github.com/tiny-sky/Tdtm/core/consts"
	"github.com/tiny-sky/Tdtm/core/dao/entity"
	"github.com/tiny-sky/Tdtm/core/dao/mysqlx/query"
	"github.com/tiny-sky/Tdtm/tools"
)

type BranchImpl struct {
	query *query.Query
}

func NewBranchImpl() BranchImpl {
	return BranchImpl{query: &query.Query{}}
}

func (g BranchImpl) CreateBatches(ctx context.Context, list entity.BranchList) error {
	err := g.query.Branch.WithContext(ctx).CreateInBatches(list, len(list))
	err = tools.WrapDbErr(err)
	return err
}

func (g BranchImpl) GetBranches(ctx context.Context, gid string) (list entity.BranchList, err error) {
	q := g.query.Branch
	list, err = g.query.Branch.WithContext(ctx).Where(q.GID.Eq(gid)).Find()
	if err = tools.WrapDbErr(err); err != nil {
		return
	}
	return
}

func (g BranchImpl) UpdateBranchStateByGid(ctx context.Context, branchId string, state consts.BranchState, errmsg string) (int64, error) {
	branch := g.query.Branch
	result, err := g.query.Branch.WithContext(ctx).
		Where(branch.BranchId.Eq(branchId)).
		UpdateSimple(branch.State.Value(string(state)), branch.LastErrMsg.Value(errmsg))
	err = tools.WrapDbErr(err)
	return result.RowsAffected, err
}
