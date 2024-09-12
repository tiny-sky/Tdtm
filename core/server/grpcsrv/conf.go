package grpcsrv

type Grpc struct {
	ListenOn string  `yaml:"listenOn"`
	Gateway  Gateway `yaml:"gateway"`
}
type Gateway struct {
	IsOpen     bool   `yaml:"isOpen"`
	CertFile   string `yaml:"certFile"`
	ServerName string `yaml:"serverName"`
}

func (grpc *Grpc) OpenGateway() bool {
	return grpc.Gateway.IsOpen
}
