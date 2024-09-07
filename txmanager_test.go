package Tdtm

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

type mockTXStore struct {
	mutex sync.Mutex
	txs   map[string]*Transaction
}

func newMockTXStore() TXStore {
	return &mockTXStore{
		txs: make(map[string]*Transaction),
	}
}

func (m *mockTXStore) CreateTX(ctx context.Context, components ...TccComponent) (string, error) {
	txid := uuid.NewString()
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.txs[txid]; ok {
		return "", fmt.Errorf("repeat txid: %s", txid)
	}

	componentTryEntities := make([]*ComponentTryEntity, 0, len(components))
	for _, component := range components {
		componentTryEntities = append(componentTryEntities, &ComponentTryEntity{
			ComponentID: component.ID(),
			TryStatus:   TryHanging,
		})
	}

	m.txs[txid] = &Transaction{
		TXID:       txid,
		Status:     TXHanging,
		CreatedAt:  time.Now(),
		Components: componentTryEntities,
	}

	return txid, nil
}

func (m *mockTXStore) TXUpdate(ctx context.Context, txID string, componentID string, accept bool) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tx, ok := m.txs[txID]
	if !ok {
		return fmt.Errorf("[TXUpdate]invalid txid: %s", txID)
	}

	for _, component := range tx.Components {
		if component.ComponentID != componentID {
			continue
		}
		if component.TryStatus != TryHanging {
			return fmt.Errorf("invalid component status: %s, componentid: %s, txid: %s", component.TryStatus, componentID, txID)
		}
		if accept {
			component.TryStatus = TrySucceesful
		} else {
			component.TryStatus = TryFailure
		}
		return nil
	}
	return fmt.Errorf("[TXUpdate]invalid component id: %s for txid: %s", componentID, txID)
}

func (m *mockTXStore) TXSubmit(ctx context.Context, txID string, success bool) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tx, ok := m.txs[txID]
	if !ok {
		return fmt.Errorf("[TXSubmit]invalid txid: %s", txID)
	}

	if success {
		if tx.Status != TXHanging && tx.Status != TXSuccessful {
			return fmt.Errorf("invalid txstatus: %s, txid: %s", tx.Status, txID)
		}
		tx.Status = TXSuccessful
	} else {
		if tx.Status != TXHanging && tx.Status != TXFailure {
			return fmt.Errorf("invalid txstatus: %s, txid: %s", tx.Status, txID)
		}
		tx.Status = TXFailure
	}
	return nil
}

func (m *mockTXStore) GetHangingTXs(ctx context.Context) ([]*Transaction, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var hangingTXs []*Transaction
	for _, tx := range m.txs {
		if tx.Status != TXHanging {
			continue
		}
		hangingTXs = append(hangingTXs, tx)
	}
	return hangingTXs, nil
}

func (m *mockTXStore) GetTX(ctx context.Context, txID string) (*Transaction, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tx, ok := m.txs[txID]
	if !ok {
		return nil, fmt.Errorf("[GetTX]invalid txid: %s", txID)
	}
	return tx, nil
}

func (m *mockTXStore) Lock(ctx context.Context, expireDuration time.Duration) error {
	// m.mutex.Lock()
	return nil
}

func (m *mockTXStore) Unlock(ctx context.Context) error {
	// m.mutex.Unlock()
	return nil
}

type Status string

const (
	StatusTried     = "tried"
	StatusConfirmed = "confirmed"
	StatusCanceled  = "canceled"
)

type mockComponent struct {
	id            string
	mutex         sync.Mutex
	statusMachine map[string]Status
}

func newMockComponent(id string) TccComponent {
	return &mockComponent{
		id:            id,
		statusMachine: make(map[string]Status),
	}
}

func (m *mockComponent) ID() string {
	return m.id
}

// 执行第一阶段的 try 操作
func (m *mockComponent) Try(ctx context.Context, req *TccReq) (*TccResp, error) {
	resp := TccResp{
		ComponentID: m.id,
		TXID:        req.TXID,
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.statusMachine[req.TXID] == StatusCanceled {
		return &resp, nil
	}

	if req.Data["reject_flag"] == true {
		m.statusMachine[req.TXID] = StatusCanceled
		return &resp, nil
	}

	if req.Data["hanging_flag"] == true {
		<-time.After(time.Second)
		return &resp, nil
	}

	if m.statusMachine[req.TXID] != StatusConfirmed {
		m.statusMachine[req.TXID] = StatusTried
	}

	resp.ACK = true
	return &resp, nil
}

func (m *mockComponent) Confirm(ctx context.Context, txID string) (*TccResp, error) {
	resp := TccResp{
		ComponentID: m.id,
		TXID:        txID,
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.statusMachine[txID] != StatusTried && m.statusMachine[txID] != StatusConfirmed {
		return &resp, nil
	}

	resp.ACK = true
	m.statusMachine[txID] = StatusConfirmed
	return &resp, nil
}

func (m *mockComponent) Cancel(ctx context.Context, txID string) (*TccResp, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.statusMachine[txID] == StatusConfirmed {
		return nil, errors.New("invalid status machine: [confirmed] when canceling")
	}

	m.statusMachine[txID] = StatusCanceled
	return &TccResp{
		ComponentID: m.id,
		ACK:         true,
		TXID:        txID,
	}, nil
}

func Test_txmanager_transaction_success(t *testing.T) {
	txmanager := NewTXManager(newMockTXStore())
	defer txmanager.stop()

	componentsCnt := 5
	componentReqs := make([]*RequestEntity, 0, componentsCnt)
	ctx := context.Background()
	for i := 0; i < componentsCnt; i++ {
		componentID := cast.ToString(i)
		if err := txmanager.Register(newMockComponent(componentID)); err != nil {
			t.Error()
			return
		}
		componentReqs = append(componentReqs, &RequestEntity{
			ComponentID: componentID,
		})
	}

	txid, ok, err := txmanager.Transaction(ctx, componentReqs...)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, true, ok)
	tx, err := txmanager.txStore.GetTX(ctx, txid)
	if err != nil {
		t.Error()
		return
	}
	assert.Equal(t, TXSuccessful, tx.Status)
}

func Test_txmanager_transaction_fail(t *testing.T) {
	txmanager := NewTXManager(newMockTXStore())
	defer txmanager.Stop()

	// 注册 5 个 component
	componentsCnt := 5
	componentReqs := make([]*RequestEntity, 0, componentsCnt)
	ctx := context.Background()
	for i := 0; i < componentsCnt; i++ {
		componentID := cast.ToString(i)
		if err := txmanager.Register(newMockComponent(componentID)); err != nil {
			t.Error(err)
			return
		}
		componentReqs = append(componentReqs, &RequestEntity{
			ComponentID: componentID,
			Request: map[string]interface{}{
				"reject_flag": true,
			},
		})
	}

	txid, ok, err := txmanager.Transaction(ctx, componentReqs...)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, false, ok)
	tx, err := txmanager.txStore.GetTX(ctx, txid)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, TXFailure, tx.Status)
}

func Test_txmanager_transaction_advance_progress(t *testing.T) {
	txmanager := NewTXManager(newMockTXStore(), WithMonitorTick(100*time.Millisecond))
	defer txmanager.stop()

	// 注册 5 个 component
	componentsCnt := 5
	componentReqs := make([]*RequestEntity, 0, componentsCnt)
	ctx := context.Background()
	for i := 0; i < componentsCnt; i++ {
		componentID := cast.ToString(i)
		if err := txmanager.Register(newMockComponent(componentID)); err != nil {
			t.Error(err)
			return
		}
		componentReqs = append(componentReqs, &RequestEntity{
			ComponentID: componentID,
			Request: map[string]interface{}{
				"hanging_flag": true,
			},
		})
	}

	txid, ok, err := txmanager.Transaction(ctx, componentReqs...)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, false, ok)
	tx, err := txmanager.txStore.GetTX(ctx, txid)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, TXFailure, tx.Status)
}

func Test_txmanager_transaction_concurrent(t *testing.T) {
	txmanager := NewTXManager(newMockTXStore(), WithMonitorTick(0), WithTimeout(0))
	defer txmanager.Stop()

	// 注册 10 个 component
	componentsCnt := 10
	for i := 0; i < componentsCnt; i++ {
		componentID := cast.ToString(i)
		if err := txmanager.Register(newMockComponent(componentID)); err != nil {
			t.Error(err)
			return
		}
	}

	// 并发 100 个分布式事务，随机取 3 个 component
	ctx := context.Background()
	concurrentTXs := 100
	componentReqCnt := 3
	var wg sync.WaitGroup
	for i := 0; i < concurrentTXs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rander := rand.New(rand.NewSource(time.Now().UnixNano()))
			componentSet := make(map[string]struct{}, componentReqCnt)
			for len(componentSet) < componentReqCnt {
				componentID := cast.ToString(rander.Intn(componentsCnt))
				componentSet[componentID] = struct{}{}
			}

			componentReqs := make([]*RequestEntity, 0, componentReqCnt)
			for componentID := range componentSet {
				componentReqs = append(componentReqs, &RequestEntity{
					ComponentID: componentID,
				})
			}

			txid, ok, err := txmanager.Transaction(ctx, componentReqs...)
			if err != nil {
				t.Error(err)
				return
			}
			assert.Equal(t, true, ok)
			tx, err := txmanager.txStore.GetTX(ctx, txid)
			if err != nil {
				t.Error(err)
				return
			}
			assert.Equal(t, TXSuccessful, tx.Status)
		}()
	}

	wg.Wait()
}
