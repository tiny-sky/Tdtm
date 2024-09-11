package client

import "github.com/tiny-sky/Tdtm/core/consts"

type GroupOpt func(group *Group)

type Group struct {
	tranType   consts.TransactionType
	levelFixed bool
	branches   []*Branch
}

func NewTccGroup(tryUri, confirmUri, cancelUri string, opts ...GroupOpt) *Group {
	g := &Group{
		tranType: consts.TCC,
	}
	for _, opt := range opts {
		opt(g)
	}

	g.branches = []*Branch{
		NewBranch(tryUri, consts.Try),
		NewBranch(confirmUri, consts.Confirm),
		NewBranch(cancelUri, consts.Cancel),
	}
	return g
}

func NewSagaGroup(normalUri, compensation string, opts ...GroupOpt) *Group {
	g := &Group{
		tranType: consts.SAGA,
	}
	for _, opt := range opts {
		opt(g)
	}

	g.branches = []*Branch{
		NewBranch(normalUri, consts.Normal),
		NewBranch(compensation, consts.Compensation),
	}
	return g
}

func (g *Group) GetTranType() consts.TransactionType {
	return g.tranType
}

func (g *Group) SetData(data []byte) *Group {
	g.set(func(branch *Branch) {
		branch.SetData(data)
	})
	return g
}

func (g *Group) SetTimeout(second int) *Group {
	g.set(func(branch *Branch) {
		branch.timeout = int64(second)
	})
	return g
}

func (g *Group) SetHeader(data []byte) *Group {
	g.set(func(branch *Branch) {
		branch.SetHeader(data)
	})
	return g
}

func (g *Group) SetLevel(level consts.Level) *Group {
	if g.levelFixed {
		return g
	}
	g.set(func(branch *Branch) {
		branch.SetLevel(level)
	})
	return g
}

func (g *Group) set(fn func(branch *Branch)) {
	for _, branch := range g.branches {
		fn(branch)
	}
}
