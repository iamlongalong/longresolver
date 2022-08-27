package k8sresolver

import (
	"github.com/iamlongalong/longresolver"
	"google.golang.org/grpc/resolver"
)

type nopResolver struct {
	cc resolver.ClientConn
}

func (r *nopResolver) Close() {
	longresolver.GetDfLogger().Infof("resolver closed")
}

// ResolveNow 没必要
func (r *nopResolver) ResolveNow(options resolver.ResolveNowOptions) {
}
