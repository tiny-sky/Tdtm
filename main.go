package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/tiny-sky/Tdtm/conf/file"
	"github.com/tiny-sky/Tdtm/core"
	"github.com/tiny-sky/Tdtm/core/coordinator"
	"github.com/tiny-sky/Tdtm/core/coordinator/executor"
	"github.com/tiny-sky/Tdtm/core/dao"
	"github.com/tiny-sky/Tdtm/core/server/runner"
	"github.com/tiny-sky/Tdtm/log"
)

var filepath = flag.String("f", "/conf.yml", "configuration file")

func main() {

	flag.Parse()
	c := file.NewFile(*filepath)

	settings, err := c.Load()
	if err != nil {
		log.Fatalf("%s", err)
	}

	settings.Init()

	dao := dao.GetTransaction()
	newCoordinator := coordinator.NewCoordinator(dao, executor.NewExecutor(), settings.AutomaticExecution2)

	// TODO : Added Other Server
	var servers []core.Server

	cronServer := runner.New(newCoordinator, dao, runner.WithMaxTimes(settings.Cron.MaxTimes), runner.WithTimeInterval(settings.Cron.TimeInterval))
	servers = append(servers, cronServer)

	var opts []core.Option
	opts = append(opts, core.WithServers(servers...))

	if settings.SetRegistry() {
		registry, err := settings.GetRegistry()
		if err != nil {
			log.Fatalf("%+v", err)
		}
		opts = append(opts, core.WithRegistry(registry))
	}

	newCore := core.New(opts...)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := newCore.Run(context.Background()); err != nil {
		log.Fatalf("%+v", err)
	}
	log.Infof("easycar server is stopped")

}
