package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "Ilya"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	//name = flag.String("name", defaultName, "Name of struct and needed field")
)

func main() {
	//flag.Parse()
	var name string
	fmt.Fscan(os.Stdin, &name)
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("%s", r.GetMessage())

	error := http.ListenAndServe("localhost:50050", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if error != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
