package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"github.com/glerchundi/grpc-helloworld/internal/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func sayHello(client helloworld.GreeterClient, req *helloworld.HelloRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rep, err := client.SayHello(ctx, req)
	if err != nil {
		log.Fatalf("%v.SayHello(_) = _, %v: ", client, err)
	}
	log.Println(rep)
}

func sayRepetitiveHello(client helloworld.GreeterClient, req *helloworld.HelloRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	stream, err := client.SayRepetitiveHello(ctx, req)
	if err != nil {
		log.Fatalf("%v.SayRepetitiveHello(_) = _, %v", client, err)
	}
	for {
		rep, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.SayRepetitiveHello(_) = _, %v", client, err)
		}
		log.Println(rep)
	}
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterClient(conn)

	sayHello(client, &helloworld.HelloRequest{Name: "Gorka"})
	sayRepetitiveHello(client, &helloworld.HelloRequest{Name: "Gorka"})
}
