package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
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
		cert, err := mytls.FSByte(false, "/server.pem")
		if err != nil {
			log.Fatalf("failed to load server certificate: %v", err)
		}

		key, err := mytls.FSByte(false, "/server.key")
		if err != nil {
			log.Fatalf("failed to load server private key: %v", err)
		}

		pair, err := tls.X509KeyPair(cert, key)
		if err != nil {
			log.Fatalf("failed to create key pair: %v", err)
		}

		pool := x509.NewCertPool()
		ok := pool.AppendCertsFromPEM(cert)
		if !ok {
			log.Fatalf("failed to append cert to pool: %v", err)
		}

		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{pair},
			NextProtos:   []string{"h2"},
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
