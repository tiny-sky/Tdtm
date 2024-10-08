package entity

import (
	"time"

	"github.com/tiny-sky/Tdtm/core/consts"
)

type Global struct {
	GID          string             `gorm:"column:g_id;type:varchar(255);not null"`               // global id
	State        consts.GlobalState `gorm:"column:state;type:varchar(255);not null;default:init"` // global State
	EndTime      int64              `gorm:"column:end_time;type:int;not null;default:0"`          // end time for the transaction
	TryTimes     int64              `gorm:"column:try_times;type:int;not null;default:0"`         // try times
	NextCronTime int64              `gorm:"column:next_cron_time;type:int;not null;default:0"`    // next cron time
	CreateTime   int64              `gorm:"create_time;autoCreateTime" json:"create_time"`        // create time
	UpdateTime   int64              `gorm:"update_time;autoCreateTime" json:"update_time"`        // last update time
}

// 指定表名
func (g *Global) TableName() string {
	return "global"
}

func NewGlobal(gId string) *Global {
	return &Global{
		GID: gId,
		// nextCronTime in first time.
		NextCronTime: time.Now().Add(3 * time.Minute).Unix(),
	}
}

func (g *Global) Phase1() bool {
	return g.State == consts.Phase1Preparing
}

func (g *Global) Phase2() bool {
	return g.GotoRollback() || g.GotoCommit()
}

func (g *Global) GotoCommit() bool {
	return g.State == consts.Phase1Success ||
		g.State == consts.Phase2Committing ||
		g.State == consts.Phase2CommitFailed // retry
}

func (g *Global) GotoRollback() bool {
	return g.State == consts.Phase1Failed ||
		g.State == consts.Phase2Rollbacking ||
		g.State == consts.Phase2RollbackFailed
}

func (g *Global) Init() bool {
	return g.State == consts.Init
}

func (g *Global) IsEmpty() bool {
	return g.GID == ""
}

func (g *Global) SetGId(gId string) {
	g.GID = gId
}

func (g *Global) GetGId() string {
	return g.GID
}
func (g *Global) SetState(state consts.GlobalState) {
	g.State = state
}

func (g *Global) GetState() consts.GlobalState {
	return g.State
}

func (g *Global) GetEndTime() int64 {
	return g.EndTime
}

func (g *Global) AllowRegister() bool {
	return g.Init()
}
