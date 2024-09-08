package conf

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/tiny-sky/Tdtm/dao/mongodb"
	"github.com/tiny-sky/Tdtm/dao/mysqlx"
	"github.com/tiny-sky/Tdtm/registry"
	"github.com/tiny-sky/Tdtm/registry/consulx"
	"github.com/tiny-sky/Tdtm/registry/etcdx"
	"github.com/tiny-sky/Tdtm/server/httpsrv"
	"github.com/tiny-sky/Tdtm/transport/common"
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
		Driver  string           `yaml:"driver"`
		Mysql   mysqlx.Settings  `yaml:"mysql"`
		Mongodb mongodb.Settings `yaml:"mongodb"`
	}

	Server struct {
		Http httpsrv.Http `yaml:"http"`
		// Grpc grpcsrv.Grpc `yaml:"grpc"`
	}

	RegistrySettings struct {
		Etcd   etcdx.Conf   `yaml:"etcd"`
		Consul consulx.Conf `yaml:"consul"`
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
	case "mongodb":
		db.Mongodb.Init()
	default:
		panic(fmt.Errorf("no support %s database", db.Driver))
	}
}

// TODO : Added tracing
func (s *Settings) Init() {
	s.DB.Init()

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	if s.Http.ListenOn == "" {
		s.Http.ListenOn = "8087"
	}
	if s.Timeout > 0 {
		common.ReplaceTimeout(time.Duration(s.Timeout) * time.Second)
	}
}

func (s *Settings) SetRegistry() bool {
	return !s.Registry.Etcd.Empty() || !s.Registry.Consul.Empty()
}

func (s *Settings) GetRegistry() (registry.Registry, error) {
	if !s.Registry.Etcd.Empty() {
		return etcdx.New(s.Registry.Etcd)
	}

	// consul and add others?
	client, err := api.NewClient(s.Registry.Consul.Conf())
	if err != nil {
		return nil, err
	}
	return consulx.New(client), nil
}
