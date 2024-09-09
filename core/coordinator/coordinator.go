package coordinator

import (
	"context"
	"fmt"
	"time"

	"github.com/tiny-sky/Tdtm/core/consts"
	"github.com/tiny-sky/Tdtm/core/dao"
	"github.com/tiny-sky/Tdtm/core/entity"
	"github.com/tiny-sky/Tdtm/log"
	"github.com/tiny-sky/Tdtm/tools"
)

var (
	ErrAutomatic = fmt.Errorf("transaction phase2 has been processed automatically")
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

func (c *Coordinator) Begin(ctx context.Context) (string, error) {
	gid := entity.GitGid()
	g := entity.NewGlobal(gid)
	g.SetState(consts.Init)
	now := time.Now().Unix()
	g.CreateTime = now
	g.UpdateTime = now
	err := c.dao.CreateGlobal(ctx, g)
	return gid, err
}

func (c *Coordinator) Close(ctx context.Context) error {
	return c.closeFn(ctx)
}

func (c *Coordinator) Register(ctx context.Context, branches entity.BranchList) error {
	if len(branches) == 0 {
		return nil
	}
	return c.dao.CreateBatches(ctx, branches)
}

func (c *Coordinator) Start(ctx context.Context, global *entity.Global) error {
	if err := c.Phase1(ctx, global); err != nil {
		return err
	}
	if c.automaticExecution2 {
		log.Infof("[Coordinator] Phase2 start gid:%v", global.GID)
		tools.GoSafe(func() {
			if err := c.Phase2(context.Background(), global); err != nil {
				log.Errorf("[Start] Phase2:err:%v", err)
				return
			}
		})
	}
	return nil
}

func (c *Coordinator) Phase1(ctx context.Context, global *entity.Global) (err error) {
	phase1State := consts.Phase1Success
	defer func() {
		if err != nil {
			phase1State = consts.Phase1Failed
		}
		global.State = phase1State
		log.Infof("[Coordinator] phase1 end gid:%v state:%v", global.GetGId(), global.State)
		if erro := c.UpdateGlobalState(ctx, global.GetGId(), phase1State); erro != nil {
			log.Errorf("[Coordinator]Phase1 UpdateGlobalState:%v", erro)
		}
	}()
	err = c.executor.Phase1(ctx, global)
	return
}

func (c *Coordinator) Phase2(ctx context.Context, global *entity.Global) (err error) {
	isRollback := global.GotoRollback()

	var (
		processingStateVal, overStateVal interface{}
	)
	processingStateVal = tools.IF(isRollback, consts.Phase2Rollbacking, consts.Phase2Committing)
	if _, err = c.dao.UpdateGlobalStateByGid(ctx, global.GetGId(), processingStateVal.(consts.GlobalState)); err != nil {
		return
	}

	overStateVal = tools.IF(isRollback, consts.Rollbacked, consts.Committed)

	defer func() {
		if err != nil {
			overStateVal = tools.IF(isRollback, consts.Phase2RollbackFailed, consts.Phase2CommitFailed)
			log.Infof("[Coordinator] Phase2 end gid %v,state:%v", global.GID, overStateVal)
			if _, erro := c.dao.UpdateGlobalStateByGid(ctx, global.GetGId(), overStateVal.(consts.GlobalState)); erro != nil {
				log.Errorf("[Phase2]UpdateGlobalStateByGid gid:%v err:%v", global.GetGId(), erro)
			}
		}
	}()
	err = c.executor.Phase2(ctx, global)
	return
}

func (c *Coordinator) Commit(ctx context.Context, global *entity.Global) error {
	if c.automaticExecution2 {
		log.Warnf("[Commit] gid:%v,warn:%v", global.GetGId(), ErrAutomatic)
		return nil
	}
	return c.Phase2(ctx, global)
}

func (c *Coordinator) Rollback(ctx context.Context, global *entity.Global) error {
	if c.automaticExecution2 {
		log.Warnf("[Rollback] gid:%v,warn:%v", global.GetGId(), ErrAutomatic)
		return nil
	}
	return c.Phase2(ctx, global)
}

func (c *Coordinator) GetGlobal(ctx context.Context, gid string) (entity.Global, error) {
	return c.dao.GetGlobal(ctx, gid)
}

func (c *Coordinator) GetBranchList(ctx context.Context, gid string) (list []*entity.Branch, err error) {
	return c.dao.GetBranches(ctx, gid)
}

func (c *Coordinator) UpdateGlobalState(ctx context.Context, gid string, state consts.GlobalState) error {
	_, err := c.dao.UpdateGlobalStateByGid(ctx, gid, state)
	return err
}
