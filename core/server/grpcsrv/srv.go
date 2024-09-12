package grpcsrv

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tiny-sky/Tdtm/core/coordinator"
	"github.com/tiny-sky/Tdtm/core/server/httpsrv"
	"github.com/tiny-sky/Tdtm/log"
	"github.com/tiny-sky/Tdtm/proto"
	"github.com/tiny-sky/Tdtm/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	_              proto.TdtmServer = (*GrpcSrv)(nil)
	ErrGidNotExist                  = errors.New("gid is not exits")
)

type GrpcSrv struct {
	proto.UnimplementedTdtmServer
	coordinator *coordinator.Coordinator

	lis net.Listener

	timeout    time.Duration
	listenOn   string
	once       sync.Once
	groupOpts  []grpc.ServerOption
	grpcServer *grpc.Server
}

func New(settings Grpc, coordinator *coordinator.Coordinator) (*GrpcSrv, error) {
	srv := &GrpcSrv{
		coordinator: coordinator,
		timeout:     10 * time.Second,
		listenOn:    tools.FigureOutListen(settings.ListenOn),
		once:        sync.Once{},
	}

	var err error

	srv.groupOpts = append(srv.groupOpts, grpc.MaxRecvMsgSize(5*1024*1024))
	srv.lis, err = net.Listen("tcp", srv.listenOn)
	return srv, err
}

func (s *GrpcSrv) Run(ctx context.Context) error {
	s.grpcServer = grpc.NewServer(s.groupOpts...)

	proto.RegisterTdtmServer(s.grpcServer, s)

	reflection.Register(s.grpcServer)
	go func() {
		if err := s.grpcServer.Serve(s.lis); err != nil {
			log.Fatalf("%+v", err)
		}
	}()
	log.Infof("[Grpc] grpc listen:%s", s.listenOn)
	return nil
}

func (s *GrpcSrv) Stop(ctx context.Context) (err error) {
	if s.grpcServer == nil {
		return
	}
	s.once.Do(func() {
		s.grpcServer.GracefulStop()
	})
	if err = s.coordinator.Close(ctx); err != nil {
		return
	}
	log.Infof("[GrpcSrv] stopped")
	return
}

func (s *GrpcSrv) Handler(certFile, name string) httpsrv.Handler {
	return func(ctx context.Context) (http.Handler, error) {
		var (
			err  error
			opts []grpc.DialOption
		)

		creds := insecure.NewCredentials()
		if certFile != "" {
			creds, err = credentials.NewClientTLSFromFile(certFile, name)
			if err != nil {
				return nil, err
			}
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))

		clientConn, err := grpc.DialContext(ctx, s.listenOn, opts...)
		if err != nil {
			fmt.Println(color.HiRedString("grpc DialContext:err:%v", err))
			return nil, err
		}

		mux := runtime.NewServeMux()
		err = proto.RegisterTdtmHandler(ctx, mux, clientConn)
		return mux, err
	}
}

func (s *GrpcSrv) Endpoin() *url.URL {
	return &url.URL{
		Scheme: "grpc",
		Host:   s.listenOn,
	}
}
