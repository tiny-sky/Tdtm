package executor

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/tiny-sky/Tdtm/core/consts"
	"github.com/tiny-sky/Tdtm/core/dao"
	"github.com/tiny-sky/Tdtm/core/dao/entity"
	"github.com/tiny-sky/Tdtm/core/transport"
	"github.com/tiny-sky/Tdtm/core/transport/common"
	"github.com/tiny-sky/Tdtm/log"
	"github.com/tiny-sky/Tdtm/tools/retry"
	"golang.org/x/sync/errgroup"
)

var DefaultExecutor = &Default{}

type (
	FilterFn func(branch *entity.Branch) bool

	Default struct {
		manager transport.Manager
		timeout time.Duration
	}
)

func NewExecutor() *Default {
	return DefaultExecutor
}

func (e *Default) Close(ctx context.Context) error {
	return e.manager.Close(ctx)
}

func (e *Default) Phase1(ctx context.Context, g *entity.Global) error {
	branches, err := dao.GetTransaction().GetBranches(ctx, g.GetGId())
	if err != nil {
		return err
	}
	return e.execute(ctx, true, branches, func(branch *entity.Branch) bool {
		if branch.Success() {
			return false
		}

		return branch.TccTry() || branch.SAGANormal()
	})
}

func (e *Default) Phase2(ctx context.Context, g *entity.Global) error {
	branches, err := dao.GetTransaction().GetBranches(ctx, g.GetGId())
	if err != nil {
		return err
	}

	return e.execute(ctx, false, branches, func(branch *entity.Branch) bool {
		if branch.Success() {
			return false
		}

		if g.GotoCommit() {
			return branch.TccConfirm()
		}

		if branch.SAGACompensation() || branch.TccCancel() {
			return true
		}
		return false
	})
}

func (e *Default) stratify(branches entity.BranchList) []entity.BranchList {
	layerList := make([]entity.BranchList, len(branches))

	sort.Slice(branches, func(i, j int) bool {
		return branches[i].Level < branches[j].Level
	})

	var (
		previousLevel consts.Level = 1 // first level default 1
		bucketIndex                = 0 // first level bucket index default 0
	)

	for i, branch := range branches {
		if i == 0 {
			previousLevel = branch.Level
		}
		if branch.Level > previousLevel {
			bucketIndex += 1
			previousLevel = branch.Level
		}
		layerList[bucketIndex] = append(layerList[bucketIndex], branch)
	}
	return layerList
}

func (e *Default) execute(ctx context.Context, shouldStratify bool, branches entity.BranchList, filterFn FilterFn) error {

	// filter branches
	for i := 0; i < len(branches); {
		if !filterFn(branches[i]) {
			branches = append(branches[:i], branches[i+1:]...)
		} else {
			i++
		}
	}

	layeredList := []entity.BranchList{branches}
	if shouldStratify {
		layeredList = e.stratify(branches)
	}

	for _, branchGroup := range layeredList {
		if len(branchGroup) == 0 {
			continue
		}

		errGroup, _ := errgroup.WithContext(ctx)
		for _, branch := range branchGroup {
			b := branch
			errGroup.Go(func() error {
				var (
					err         error
					errmsg      string
					branchState = consts.BranchSucceed
				)

				// request the RM
				if err := e.request(ctx, b); err != nil {
					log.Errorf("[Executor] request branch: %s err:%v", b.BranchId, err)
					branchState = consts.BranchFailState
					errmsg = err.Error()
				}

				b.State = branchState
				if _, erro := dao.GetTransaction().UpdateBranchStateByGid(ctx, b.BranchId, b.State, errmsg); erro != nil {
					log.Errorf("[Executor]update branch:%s state error:%v", b.BranchId, erro)
				}
				return err
			})
		}

		//in phase1, we have to stop execution and don't go to the next level RM if some RM is wrong
		if err := errGroup.Wait(); err != nil {
			return err
		}
	}
	return nil
}

func (e *Default) request(ctx context.Context, b *entity.Branch) (err error) {
	transporter, err := e.manager.GetTransporter(common.Net(b.Protocol))
	if err != nil {
		return fmt.Errorf("[Executor]branchid:%s get transport error:%v", b.BranchId, err)
	}

	defer func() {
		if err != nil {
			if errors.Is(err, retry.ErrOverMaximumAttempt) {
				log.Warnf("over maximum attempt")
			}
			err = fmt.Errorf("[Executor] Request branchid:%s request error:%v", b.BranchId, err)
		}
	}()

	var reqOpts []common.Option

	timeout := e.timeout
	if b.Timeout > 0 {
		timeout = time.Second * time.Duration(b.Timeout)
	}
	reqOpts = append(reqOpts, common.WithTimeout(timeout))

	req := common.NewReq([]byte(b.ReqData), []byte(b.ReqHeader), reqOpts...)
	req.AddHeaders(b.GID, b.BranchId)

	r := retry.New(2, retry.WithMaxBackOffTime(1*time.Second))

	err = r.Run(func() error {
		_, err = transporter.Request(ctx, b.Url, req)
		return err
	})
	return
}
