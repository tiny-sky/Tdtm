package grpcsrv

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/tiny-sky/Tdtm/core/coordinator"
	"github.com/tiny-sky/Tdtm/proto"
	"github.com/tiny-sky/Tdtm/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	_             proto.TdtmServer = (*GrpcSrv)(nil)
	ErrGidNoExist                  = errors.New("gid is not exits")
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

	if settings.Tls() {
		certificate, err := tls.LoadX509KeyPair(settings.CertFile, settings.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		tlsConf := &tls.Config{Certificates: []tls.Certificate{certificate}}
		srv.groupOpts = append(srv.groupOpts, grpc.Creds(credentials.NewTLS(tlsConf)))
	}
	srv.groupOpts = append(srv.groupOpts, grpc.MaxRecvMsgSize(5*1024*1024))
	srv.lis, err = net.Listen("tcp", srv.listenOn)
	return srv, err
}
