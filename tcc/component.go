package tcc

import "context"

// tcc 请求参数
type TccReq struct {
	TccResp

	ComponentID string                 `json:"componentID"`
	TXID        string                 `json:"txID"`
	Data        map[string]interface{} `json:"data"`
}

// tcc 响应结果
type TccResp struct {
	ComponentID string `json:"componentID"`
	ACK         bool   `json:"ack"`
	TXID        string `json:"txID"`
}

type TccComponent interface {
	// 返回组件唯一 id
	ID() string
	// 执行第一阶段的 try 操作
	Try(ctx context.Context, req *TccReq) (*TccResp, error)
	// 执行第二阶段的 confirm 操作
	Confirm(ctx context.Context, txID string) (*TccResp, error)
	// 执行第二阶段的 cancel 操作
	Cancel(ctx context.Context, txID string) (*TccResp, error)
}
