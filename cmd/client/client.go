package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tiny-sky/Tdtm/client"
	"github.com/tiny-sky/Tdtm/log"
)

func main() {
	var opts []client.Option
	opts = append(opts, client.WithConnTimeout(5*time.Second))

	cli, err := client.New("127.0.0.1:1234", opts...)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	ctx := context.Background()
	defer cli.Close(ctx)

	gid, err := cli.Begin(ctx)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	var groups []*client.Group
	if err = cli.Register(ctx, gid, groups); err != nil {
		log.Fatalf("%+v", err)
	}

	if err := cli.Start(ctx, gid); err != nil {
		fmt.Println("start err:", err)
	}
	fmt.Printf("[%s] Transaction completed\n", gid)
}
