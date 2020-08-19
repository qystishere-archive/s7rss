package provider

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/qystishere/s7rss/storage"
)

type Provider struct {
	newsStorage storage.NewsStorager

	server *grpc.Server
}

func New(newsStorage storage.NewsStorager) *Provider {
	return &Provider{
		newsStorage: newsStorage,
	}
}

func (p *Provider) Listen(host string, port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}

	p.server = grpc.NewServer([]grpc.ServerOption{}...)
	RegisterProviderServer(p.server, p)
	return p.server.Serve(listener)
}

func (p *Provider) Stop() {
	p.server.GracefulStop()
}
