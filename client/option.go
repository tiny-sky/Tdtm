package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// DefaultOptions default for a Client Options
var DefaultOptions = &Options{
	connTimeout: 15 * time.Second,
	dailOpts:    []grpc.DialOption{},
}

type HandlerFn func(ctx context.Context) error

type Option func(options *Options)

type Options struct {
	connTimeout time.Duration
	dailOpts    []grpc.DialOption
	isDiscovery bool
	beforeFunc  HandlerFn
	afterFunc   HandlerFn
}

func WithConnTimeout(seconds time.Duration) Option {
	return func(options *Options) {
		if seconds > 0 {
			options.connTimeout = seconds
		}
	}
}

func WithBeforeFunc(before HandlerFn) Option {
	return func(options *Options) {
		options.beforeFunc = before
	}
}

func WithAfterFunc(after HandlerFn) Option {
	return func(options *Options) {
		options.afterFunc = after
	}
}

func WithDiscovery() Option {
	return func(options *Options) {
		options.isDiscovery = true
	}
}

func WithGrpcDailOpts(opts []grpc.DialOption) Option {
	return func(options *Options) {
		options.dailOpts = append(options.dailOpts, opts...)
	}
}
