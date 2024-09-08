package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/tiny-sky/Tdtm/conf/file"
	"github.com/tiny-sky/Tdtm/dao"
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
	newCoordinator := coordinator
}
