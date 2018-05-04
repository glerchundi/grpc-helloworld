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
	addr      = flag.String("addr", "", "")
	authority = flag.String("authority", "", "")
)

func sayHello(client helloworld.GreeterClient, req *helloworld.SayHelloRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rep, err := client.SayHello(ctx, req)
	if err != nil {
		log.Fatalf("%v.SayHello(_) = _, %v: ", client, err)
	}
	log.Println(rep)
}

func sayRepetitiveHello(client helloworld.GreeterClient, req *helloworld.SayRepetitiveHelloRequest) {
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

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(
			credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			}),
		),
	}

	conn, err := grpc.Dial(*addr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := helloworld.NewGreeterClient(conn)

	sayHello(client, &helloworld.SayHelloRequest{Name: "Gorka"})
	sayRepetitiveHello(client, &helloworld.SayRepetitiveHelloRequest{Name: "Gorka", Count: 3})
}
