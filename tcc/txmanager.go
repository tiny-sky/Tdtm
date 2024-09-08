package tcc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/tiny-sky/Tdtm/log"
)

type TXManager struct {
	ctx            context.Context
	stop           context.CancelFunc
	opts           *Options
	txStore        TXStore
	registryCenter *registryCenter
}

func NewTXManager(txStore TXStore, opts ...Option) *TXManager {
	ctx, cancel := context.WithCancel(context.Background())
	txManager := TXManager{
		opts:           &Options{},
		txStore:        txStore,
		registryCenter: newRegistryCenter(),
		ctx:            ctx,
		stop:           cancel,
	}

	for _, opt := range opts {
		opt(txManager.opts)
	}

	repair(txManager.opts)

	go txManager.run()
	return &txManager
}

func (t *TXManager) backOffTick(tick time.Duration) time.Duration {
	tick <<= 1
	if threshold := t.opts.MonitorTick << 3; tick > threshold {
		return threshold
	}
	return tick
}

func (t *TXManager) Stop() {
	t.stop()
}

func (t *TXManager) Register(component TccComponent) error {
	return t.registryCenter.register(component)
}

// 事务
func (t *TXManager) Transaction(ctx context.Context, reqs ...*RequestEntity) (string, bool, error) {
	tctx, cancel := context.WithTimeout(ctx, t.opts.Timeout)
	defer cancel()

	// 获得所有的组件
	componentEntities, err := t.getComponents(tctx, reqs...)
	if err != nil {
		return "", false, err
	}

	// 创建事务ID
	txID, err := t.txStore.CreateTX(tctx, componentEntities.ToComponents()...)
	if err != nil {
		return "", false, err
	}

	// 两阶段提交
	return txID, t.twoPhaseCommit(ctx, txID, componentEntities), nil
}

func (t *TXManager) run() {
	var tick time.Duration
	var err error
	for {
		// 如果出现了失败，tick 需要避让，遵循退避策略增大 tick 间隔时长
		if err == nil {
			tick = t.opts.MonitorTick
		} else {
			tick = t.backOffTick(tick)
		}

		select {
		case <-t.ctx.Done():
			return

		case <-time.After(tick):
			// 加锁，避免多个分布式多个节点的监控任务重复执行
			log.Infof("t.txStore.Lock")
			if err = t.txStore.Lock(t.ctx, t.opts.MonitorTick); err != nil {
				err = nil
				continue
			}

			// 获取仍然处于 hanging 状态的事务
			var txs []*Transaction
			if txs, err = t.txStore.GetHangingTXs(t.ctx); err != nil {
				_ = t.txStore.Unlock(t.ctx)
				continue
			}

			err = t.batchAdvanceProgress(txs)
			log.Infof("t.txStore.Unlock")
			_ = t.txStore.Unlock(t.ctx)
		}
	}
}

// 并发执行，推进各比事务的进度
func (t *TXManager) batchAdvanceProgress(txs []*Transaction) error {
	errCh := make(chan error)
	go func() {
		var wg sync.WaitGroup
		for _, tx := range txs {
			tx := tx
			wg.Add(1)

			go func() {
				defer wg.Done()

				if err := t.advanceProgress(tx); err != nil {
					// 遇到错误则投递到 errCh
					errCh <- err
				}
			}()
		}

		wg.Wait()
		close(errCh)
	}()

	var firstErr error
	for err := range errCh {
		if firstErr != nil {
			continue
		}
		firstErr = err
	}

	return firstErr
}

func (t *TXManager) advanceProgressByTXID(txID string) error {
	tx, err := t.txStore.GetTX(t.ctx, txID)
	if err != nil {
		return err
	}
	return t.advanceProgress(tx)
}

// 推进事务进度
func (t *TXManager) advanceProgress(tx *Transaction) error {
	txStatus := tx.getStatus(time.Now().Add(-t.opts.Timeout))

	if txStatus == TXHanging {
		return nil
	}

	success := txStatus == TXSuccessful
	var confirmOrCancel func(ctx context.Context, component TccComponent) (*TccResp, error)
	var txAdvanceProgress func(ctx context.Context) error

	if success {
		confirmOrCancel = func(ctx context.Context, component TccComponent) (*TccResp, error) {
			// 对 component 进行第二阶段的 confirm 操作
			return component.Confirm(ctx, tx.TXID)
		}
		txAdvanceProgress = func(ctx context.Context) error {
			return t.txStore.TXSubmit(ctx, tx.TXID, true)
		}
	} else {
		confirmOrCancel = func(ctx context.Context, component TccComponent) (*TccResp, error) {
			// 对 component 进行第二阶段的 cancel 操作
			return component.Cancel(ctx, tx.TXID)
		}

		txAdvanceProgress = func(ctx context.Context) error {
			return t.txStore.TXSubmit(ctx, tx.TXID, false)
		}
	}

	// 处理事务
	for _, component := range tx.Components {
		components, err := t.registryCenter.getComponents(component.ComponentID)
		if err != nil || len(components) == 0 {
			return errors.New("get tcc component failed")
		}
		// 执行二阶段的 confirm 或者 cancel 操作
		resp, err := confirmOrCancel(t.ctx, components[0])
		if err != nil {
			return err
		}
		if !resp.ACK {
			return fmt.Errorf("component: %s ack failed", component.ComponentID)
		}
	}

	// 提交事务
	return txAdvanceProgress(t.ctx)
}

func (t *TXManager) twoPhaseCommit(ctx context.Context, txID string, componentEntities ComponentEntities) bool {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 并发执行，只要中间某次出现了失败，直接终止流程进行 cancel
	// 如果全量执行成功，则批量执行 confirm，然后返回成功的 ack
	errCh := make(chan error, len(componentEntities))

	go func() {
		var wg sync.WaitGroup
		for _, componentEntity := range componentEntities {
			componentEntity := componentEntity
			wg.Add(1)

			go func() {
				defer wg.Done()
				resp, err := componentEntity.Component.Try(cctx, &TccReq{
					ComponentID: componentEntity.Component.ID(),
					TXID:        txID,
					Data:        componentEntity.Request,
				})

				if err != nil || !resp.ACK {
					log.ErrorContextf(cctx, "tx try failed, tx id: %s, comonent id: %s, err: %v", txID, componentEntity.Component.ID(), err)

					if _err := t.txStore.TXUpdate(cctx, txID, componentEntity.Component.ID(), false); _err != nil {
						log.ErrorContextf(cctx, "tx updated failed, tx id: %s, component id: %s, err: %v", txID, componentEntity.Component.ID(), _err)
					}
					errCh <- fmt.Errorf("component: %s try failed", componentEntity.Component.ID())
					return
				}

				if err = t.txStore.TXUpdate(cctx, txID, componentEntity.Component.ID(), true); err != nil {
					log.ErrorContextf(cctx, "tx updated failed, tx id: %s, component id: %s, err: %v", txID, componentEntity.Component.ID(), err)
					errCh <- err
				}
			}()
		}

		wg.Wait()
		close(errCh)
	}()

	successful := true
	if err := <-errCh; err != nil {
		cancel()
		successful = false
	}

	if err := t.advanceProgressByTXID(txID); err != nil {
		log.ErrorContextf(ctx, "advance tx progress fail, txid: %s, err: %v", txID, err)
	}
	return successful
}

func (t *TXManager) getComponents(ctx context.Context, reqs ...*RequestEntity) (ComponentEntities, error) {
	if len(reqs) == 0 {
		return nil, errors.New("emtpy task")
	}

	idToReq := make(map[string]*RequestEntity, len(reqs))
	componentIDs := make([]string, 0, len(reqs))

	for _, req := range reqs {
		if _, ok := idToReq[req.ComponentID]; ok {
			return nil, fmt.Errorf("repeat component : %s", req.ComponentID)
		}
		idToReq[req.ComponentID] = req
		componentIDs = append(componentIDs, req.ComponentID)
	}

	components, err := t.registryCenter.getComponents(componentIDs...)
	if err != nil {
		return nil, err
	}
	if len(componentIDs) != len(components) {
		return nil, errors.New("invalid componentIDs ")
	}

	entities := make(ComponentEntities, 0, len(components))
	for _, component := range components {
		entities = append(entities, &ComponentEntity{
			Request:   idToReq[component.ID()].Request,
			Component: component,
		})
	}

	return entities, nil
}
