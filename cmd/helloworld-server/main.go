package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/glerchundi/grpc-helloworld/internal/helloworld"
	"github.com/glerchundi/grpc-helloworld/internal/web"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	tls      = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS key file")
	port     = flag.Int("port", 10000, "The server port")
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

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	m := cmux.New(lis)
	grpcL := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	otherL := m.Match(cmux.Any())

	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpcS := grpc.NewServer(opts...)
	helloworld.RegisterGreeterServer(grpcS, &helloworldServer{})

	httpS := &http.Server{Handler: http.FileServer(web.FS(false))}

	// Use the muxed listeners for your servers.
	go grpcS.Serve(grpcL)
	go httpS.Serve(otherL)

	// Start serving!
	m.Serve()
}
