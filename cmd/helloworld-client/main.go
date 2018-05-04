package main

import (
	"context"
	"crypto/tls"
	"flag"
	"io"
	"log"
	"time"

	"github.com/glerchundi/grpc-helloworld/internal/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	usesTLS   = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	addr      = flag.String("addr", "", "")
	authority = flag.String("authority", "", "")
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
	if *usesTLS {
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	if authority != nil {
		opts = append(opts, grpc.WithAuthority(*authority))
	}

	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterClient(conn)

	sayHello(client, &helloworld.HelloRequest{Name: "Gorka"})
	sayRepetitiveHello(client, &helloworld.HelloRequest{Name: "Gorka"})
}
