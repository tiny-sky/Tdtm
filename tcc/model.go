package tcc

import (
	"time"
)

type RequestEntity struct {
	// 组件名称
	ComponentID string `json:"componentName"`
	// 组件入参
	Request map[string]interface{} `json:"request"`
}

type ComponentEntities []*ComponentEntity

type ComponentEntity struct {
	Request   map[string]interface{}
	Component TccComponent
}

func (c ComponentEntities) ToComponents() []TccComponent {
	components := make([]TccComponent, 0, len(c))
	for _, entity := range c {
		components = append(components, entity.Component)
	}
	return components
}

// 事务状态
type TXStatus string
type ComponentTryStatus string

const (
	// 事务执行中
	TXHanging TXStatus = "hanging"
	// 事务成功
	TXSuccessful TXStatus = "successful"
	// 事务失败
	TXFailure TXStatus = "failure"
)

const (
	// 事务执行中
	TryHanging ComponentTryStatus = "hanging"
	// 事务成功
	TrySucceesful ComponentTryStatus = "successful"
	// 事务失败
	TryFailure ComponentTryStatus = "failure"
)

func (t TXStatus) String() string {
	return string(t)
}

func (c ComponentTryStatus) String() string {
	return string(c)
}

type ComponentTryEntity struct {
	ComponentID string
	TryStatus   ComponentTryStatus
}

// 事务
type Transaction struct {
	TXID       string `json:"txID"`
	Components []*ComponentTryEntity
	Status     TXStatus  `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (t *Transaction) getStatus(createBefore time.Time) TXStatus {
	var Exist bool
	for _, component := range t.Components {
		if component.TryStatus == TryFailure {
			return TXFailure
		}
		Exist = Exist || (component.TryStatus != TrySucceesful)
	}

	if Exist && t.CreatedAt.Before(createBefore) {
		return TXFailure
	}

	if Exist {
		return TXHanging
	}

	return TXSuccessful
}
