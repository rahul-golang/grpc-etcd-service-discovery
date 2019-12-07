package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/coreos/etcd/client"
	rolesPb "github.com/rahul-golang/grpc-etcd-service-discovery/roles-microservice/pb"
	pb "github.com/rahul-golang/grpc-etcd-service-discovery/user-microservice/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type contextKey int

const (
	clientIDKey contextKey = iota
)

type Server struct {
	users       []*pb.User
	rolesClient rolesPb.RolesClient
}

func discover_services() string {

	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi := client.NewKeysAPI(c)

	log.Print("Getting '/role' key value")
	resp, err := kapi.Get(context.Background(), "/role", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Printf("Get is done. Metadata is %q\n", resp)
		// print value
		log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
	}
	return resp.Node.Value

}

func getRolesClient() rolesPb.RolesClient {
	conn, err := grpc.Dial(discover_services(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to start gRPC connection: %v", err)
	}
	//defer conn.Close()

	return rolesPb.NewRolesClient(conn)
}

func (s *Server) GetUser(_ context.Context, req *pb.GetUserRequest) (*pb.UserReply, error) {
	if req.UserId < 0 || req.UserId > int32(len(s.users)) {

		return nil, errors.New("invalid user")
	}
	user := s.users[req.UserId]
	roleReq := &rolesPb.GetUserRoleRequest{
		UserId: req.UserId,
	}
	rolesReply, err := s.rolesClient.GetUserRole(context.Background(), roleReq)
	if err != nil {
		return nil, err
	}

	roles := make([]*pb.Role, 0)
	for _, role := range rolesReply.Roles {
		roles = append(roles, &pb.Role{
			Id:   role.Id,
			Name: role.Name,
		})
	}
	return &pb.UserReply{
		User:  user,
		Roles: roles,
	}, nil
}

// authenticateAgent check the client credentials
func authenticateClient(ctx context.Context) (string, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		clientLogin := strings.Join(md["login"], "")
		clientPassword := strings.Join(md["password"], "")

		if clientLogin != "AAA" {
			return "", fmt.Errorf("unknown user %s", clientLogin)
		}
		if clientPassword != "BBB" {
			return "", fmt.Errorf("bad password %s", clientPassword)
		}

		log.Printf("authenticated client: %s", clientLogin)

		return "42", nil
	}
	return "", fmt.Errorf("missing credentials")
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//s, ok := info.Server.(Server)
	// if !ok {
	// 	return nil, fmt.Errorf("unable to cast server")
	// }
	clientID, err := authenticateClient(ctx)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, clientIDKey, clientID)

	return handler(ctx, req)
}

func main() {
	users := []*pb.User{
		{
			Id:    1,
			Email: "bob@example.com",
			Name:  "Bob",
		},
		{
			Id:    2,
			Email: "amy@example.com",
			Name:  "Amy",
		},
		{
			Id:    3,
			Email: "george@example.com",
			Name:  "George",
		},
		{
			Id:    4,
			Email: "lily@msys.com",
			Name:  "Lily",
		},
		{
			Id:    5,
			Email: "jacob@example.com",
			Name:  "Jacob",
		},
	}

	lis, err := net.Listen("tcp", "localhost:7000")
	if err != nil {
		log.Fatalf("failed to initializa TCP listen: %v", err)
	}
	defer lis.Close()

	certFile := "server-cert.pem"
	keyFile := "server-key.pem"

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(unaryInterceptor)}

	server := grpc.NewServer(opts...)
	roleServer := &Server{
		users:       users,
		rolesClient: getRolesClient(),
	}
	pb.RegisterUsersServer(server, roleServer)

	server.Serve(lis)
}
