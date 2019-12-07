package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/rahul-golang/grpc-etcd-service-discovery/user-microservice/pb"
)

type contextKey int

const (
	// apiVersion is version of API is provided by server
	apiVersion             = "v1"
	clientIDKey contextKey = iota
)

type Authentication struct {
	Login    string
	Password string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login":    a.Login,
		"password": a.Password,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security
func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

func main() {
	// get configuration
	address := flag.String("server", "localhost:7000", "gRPC server in format host:port")
	flag.Parse()

	// Set up a connection to the server.
	//conn, err := grpc.Dial(*address, grpc.WithInsecure())
	auth := Authentication{
		Login:    "AAA",
		Password: "BBB",
	}

	creds, err := credentials.NewClientTLSFromFile("server-cert.pem", "")
	if err != nil {
		log.Fatalf("cert load error: %s", err)
	}

	conn, err := grpc.Dial(*address, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&auth))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewUsersClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call Create
	req1 := pb.GetUserRequest{
		UserId: 1,
	}
	res1, err := c.GetUser(ctx, &req1)
	if err != nil {
		log.Fatalf("Create  TODO failed: %v", err)
	}
	log.Printf("Create result: <%+v>\n\n", res1)

}
