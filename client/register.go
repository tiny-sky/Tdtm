package client

import (
	"github.com/tiny-sky/Tdtm/core/registry"
	"github.com/tiny-sky/Tdtm/core/resolver"
)

func RegisterBuilder(discovery registry.Discovery) {
	resolver.Register(discovery)
}
