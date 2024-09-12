package client

import (
	"sync/atomic"

	"github.com/tiny-sky/Tdtm/core/consts"
)

type Manger struct {
	groups []*Group
	level  consts.Level
}

func NewManger() *Manger {
	return &Manger{
		level: 1,
	}
}

func (m *Manger) AddNextWaitGroups(groups ...*Group) *Manger {
	atomic.AddUint32((*uint32)(&m.level), 1)
	return m.addGroups(groups...)
}

func (m *Manger) AddGroups(groups ...*Group) *Manger {
	return m.addGroups(groups...)
}

func (m *Manger) addGroups(groups ...*Group) *Manger {
	for _, group := range groups {
		group.SetLevel(m.level)
		m.groups = append(m.groups, group)
	}
	return m
}

func (m *Manger) Groups() []*Group {
	return m.groups
}
