package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/glerchundi/grpc-helloworld/internal/helloworld"
	mytls "github.com/glerchundi/grpc-helloworld/internal/tls"
	"github.com/glerchundi/grpc-helloworld/internal/web"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	wantsTLS = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	addr     = flag.String("addr", "", "")
	port     = flag.Int("port", 8443, "The server port")
)

type helloworldServer struct {
}

func reply(name string) *helloworld.HelloReply {
	return &helloworld.HelloReply{Message: fmt.Sprintf("Hello %s!!!", name)}
}

func (hws *helloworldServer) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return reply(req.Name), nil
}

func (hws *helloworldServer) SayRepetitiveHello(req *helloworld.HelloRequest, stream helloworld.Greeter_SayRepetitiveHelloServer) error {
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("%s-%d", req.Name, i)
		if err := stream.Send(reply(name)); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var tlsConfig *tls.Config
	if *wantsTLS {
		var hosts []string
		if *addr != "" {
			hosts = append(hosts, *addr)
		}

		hosts = append(
			hosts,
			"localhost", fmt.Sprintf("localhost:%d", port), "127.0.0.1",
		)

		key, cert, err := mytls.GenerateCertificate(hosts)
		if err != nil {
			log.Fatalf("failed to generate certiticate: %v", err)
		}

		tlsConfig, err = mytls.NewTLSConfig(key, cert)
		if err != nil {
			log.Fatalf("failed to create tls config: %v", err)
		}
	}

	m := cmux.New(lis)
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpL := m.Match(cmux.Any())

	var opts []grpc.ServerOption
	if tlsConfig != nil {
		opts = []grpc.ServerOption{grpc.Creds(credentials.NewTLS(tlsConfig))}
	}

	grpcS := grpc.NewServer(opts...)
	helloworld.RegisterGreeterServer(grpcS, &helloworldServer{})

	httpS := &http.Server{
		Handler:   http.FileServer(web.FS(false)),
		TLSConfig: tlsConfig,
	}

	if tlsConfig != nil {
		httpL = tls.NewListener(httpL, tlsConfig)
	}

	// Use the muxed listeners for your servers.
	go grpcS.Serve(grpcL)
	go httpS.Serve(httpL)

	// Start serving!
	m.Serve()
}
