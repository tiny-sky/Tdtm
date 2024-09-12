package conf

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tiny-sky/Tdtm/core/dao/mysqlx"
	"github.com/tiny-sky/Tdtm/core/registry"
	"github.com/tiny-sky/Tdtm/core/registry/etcdx"
	"github.com/tiny-sky/Tdtm/core/server/grpcsrv"
	"github.com/tiny-sky/Tdtm/core/server/httpsrv"
	"github.com/tiny-sky/Tdtm/core/transport/common"

	_ "net/http/pprof"
)

type (
	Settings struct {
		Server              `yaml:"server"`
		DB                  DB    `yaml:"db"`
		Timeout             int64 `yaml:"timeout"`
		AutomaticExecution2 bool  `yaml:"automaticExecution2"`
		// Tracing             Tracing          `yaml:"tracing"`
		Registry RegistrySettings `yaml:"registry"`
		Cron     Cron             `yaml:"cron"`
	}

	DB struct {
		Driver string          `yaml:"driver"`
		Mysql  mysqlx.Settings `yaml:"mysql"`
	}

	Server struct {
		Http httpsrv.Http `yaml:"http"`
		Grpc grpcsrv.Grpc `yaml:"grpc"`
	}

	RegistrySettings struct {
		Etcd etcdx.Conf `yaml:"etcd"`
	}

	Cron struct {
		MaxTimes     int `yaml:"maxTimes"`
		TimeInterval int `yaml:"timeInterval"`
	}
)

func (db *DB) Init() {
	switch db.Driver {
	case "mysql":
		db.Mysql.Init()
	default:
		panic(fmt.Errorf("no support %s database", db.Driver))
	}
}

// TODO : Added tracing
func (s *Settings) Init() {
	s.DB.Init()

	// 性能分析
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	if s.Http.ListenOn == "" {
		s.Http.ListenOn = "8087"
	}
	if s.Grpc.ListenOn == "" {
		s.Grpc.ListenOn = "8089"
	}
	if s.Timeout > 0 {
		common.ReplaceTimeout(time.Duration(s.Timeout) * time.Second)
	}
}

func (s *Settings) SetRegistry() bool {
	return !s.Registry.Etcd.Empty()
}

func (s *Settings) GetRegistry() (registry.Registry, error) {
	return etcdx.New(s.Registry.Etcd)
}
