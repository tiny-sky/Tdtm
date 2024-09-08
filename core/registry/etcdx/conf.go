package etcdx

import (
	"context"
	"time"
)

// Conf is the config item with the given key on etcd.
type Conf struct {
	Hosts              []string `yaml:"hosts"`
	User               string   `yaml:"user"`
	Pass               string   `yaml:"pass"`
	InsecureSkipVerify bool     `json:""`
}

func (c *Conf) Empty() bool {
	return len(c.Hosts) == 0
}

type Options struct {
	ttl time.Duration
	ctx context.Context
}

type Option func(options *Options)

func newDefault() Options {
	return Options{ttl: 10 * time.Second, ctx: context.Background()}
}
