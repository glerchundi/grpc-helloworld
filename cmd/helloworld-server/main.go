package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/glerchundi/grpc-helloworld/internal/helloworld"
	mytls "github.com/glerchundi/grpc-helloworld/internal/tls"
	"github.com/glerchundi/grpc-helloworld/internal/web"
	"google.golang.org/grpc"
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

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
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
			"localhost", fmt.Sprintf("localhost:%d", *port), "127.0.0.1",
		)

		key, cert, err := mytls.GenerateCertificate(hosts)
		if err != nil {
			log.Fatalf("failed to generate certiticate: %v", err)
		}

		tlsConfig, err = mytls.NewTLSConfig(key, cert)
		if err != nil {
			log.Fatalf("failed to create tls config: %v", err)
		}

		lis = tls.NewListener(lis, tlsConfig)
	}

	grpcServer := grpc.NewServer()
	helloworld.RegisterGreeterServer(grpcServer, &helloworldServer{})

	server := &http.Server{
		Handler: grpcHandlerFunc(grpcServer, http.FileServer(web.FS(false))),
	}

	// Start serving!
	if err := server.Serve(lis); err != nil {
		log.Fatalf("%v", err)
	}
}
