package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tiny-sky/Tdtm/core/transport/common"
	"github.com/tiny-sky/Tdtm/tools"
)

type Transport struct{}

func NewTransporter() *Transport {
	return new(Transport)
}

func (cli *Transport) GetType() common.Net {
	return common.Http
}

func (cli *Transport) Request(ctx context.Context, uri string, req *common.Req) (resp *common.Resp, err error) {
	resp, err = cli.req(ctx, uri, req)
	return resp, err
}

func (cli *Transport) req(ctx context.Context, uri string, req *common.Req) (*common.Resp, error) {
	resp, err := tools.RestyCli.SetTimeout(req.Timeout).R().SetContext(ctx).SetHeaders(req.Headers).SetBody(req.Body).Post(uri)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("[http Transporter]http status code: %d", resp.StatusCode())
	}
	return &common.Resp{
		Body: resp.Body(),
	}, nil
}

func (cli Transport) Close(ctx context.Context) error {
	return nil
}
