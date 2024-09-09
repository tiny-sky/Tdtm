package runner

import (
	"context"
	"time"

	"github.com/tiny-sky/Tdtm/core/coordinator"
	"github.com/tiny-sky/Tdtm/core/dao"
	"github.com/tiny-sky/Tdtm/core/dao/entity"
	"github.com/tiny-sky/Tdtm/log"
	"github.com/tiny-sky/Tdtm/tools"
	"github.com/tiny-sky/Tdtm/tools/retry"
)

type Runner struct {
	ticker  *time.Ticker
	backoff *retry.Retry
	cancel  func()
	ctx     context.Context

	dao         dao.TransactionDao
	coordinator *coordinator.Coordinator

	options *Options
}

func New(coordinator *coordinator.Coordinator, dao dao.TransactionDao, opts ...Option) *Runner {
	options := Default
	for _, fn := range opts {
		fn(options)
	}
	r := &Runner{
		coordinator: coordinator,
		dao:         dao,
		options:     options,
		backoff:     retry.New(10, retry.WithMaxBackOffTime(2*time.Minute), retry.WithFactor(2)),
	}
	r.ticker = time.NewTicker(r.options.duration)

	r.ctx, r.cancel = context.WithCancel(context.Background())
	go r.loop()
	return r
}

func (r *Runner) loop() {
	defer r.ticker.Stop()
	for {
		select {
		case <-r.ticker.C:
			list, err := r.dao.FindProcessingList(r.ctx, 2, r.options.MaxTimes)
			if err != nil {
				log.Errorf("[runner] loop err:%v", err)
				continue
			}

			if len(list) == 0 {
				// backoff for ticker
				duration := r.backoff.Duration()
				if duration != r.backoff.MaxBackOffTime() {
					log.Infof("[duration] Reset ticker:%v", duration)
					r.ticker.Reset(duration)
				}
				continue
			}
			r.runJob(list)
			// reset to default value
			r.ticker.Reset(r.options.duration)
			r.backoff.Reset()

		case <-r.ctx.Done():
			return
		}
	}
}

func (r *Runner) runJob(list []*entity.Global) {
	ctx := context.Background()

	for i := 0; i < len(list); i++ {
		global := list[i]

		tools.GoSafe(func() {
			var err error

			defer func() {
				if err != nil {
					log.Errorf("[Runner] err:%v", err)
				}
			}()

			if global.Phase1() {
				err = r.coordinator.Phase1(ctx, global)
			} else if global.Phase2() {
				err = r.coordinator.Phase2(ctx, global)
			} else {
				log.Warnf("[Runner] global:%v state :%v is wrong", global.GID, global.State)
			}

			nextCronime := time.Now().Add(r.options.timeInterval).Unix()
			if err := r.dao.IncrTryTimes(ctx, global.GID, int(nextCronime)); err != nil {
				log.Warnf("[Runner] update IncrTryTimes gid :%v  err:%v", global.GID, err)
			}
		})
	}
}

func (r *Runner) Run(ctx context.Context) (_ error) { return nil }

func (r *Runner) Stop(ctx context.Context) error {
	r.cancel()
	return nil
}
