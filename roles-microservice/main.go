package main

import (
	"errors"
	"log"
	"net"
	"time"

	"go.etcd.io/etcd/client"

	pb "github.com/rahul-golang/grpc-etcd-service-discovery/roles-microservice/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

//Server **
type Server struct {
	userRoles map[int32][]*pb.Role
	roles     []*pb.Role
}

//GetRoles return all roles
func (s *Server) GetRoles(_ context.Context, _ *pb.EmptyRequest) (*pb.RolesReply, error) {
	return &pb.RolesReply{
		Roles: s.roles,
	}, nil
}

//GetUserRole **
func (s *Server) GetUserRole(_ context.Context, req *pb.GetUserRoleRequest) (*pb.UserRoleReply, error) {
	roles, ok := s.userRoles[req.UserId]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &pb.UserRoleReply{
		UserId: req.UserId,
		Roles:  roles,
	}, nil
}

func main() {

	var (
		normal = &pb.Role{
			Id:   1,
			Name: "normal",
		}
		editor = &pb.Role{
			Id:   2,
			Name: "editor",
		}
		admin = &pb.Role{
			Id:   3,
			Name: "admin",
		}
		superUser = &pb.Role{
			Id:   4,
			Name: "super user",
		}
	)

	lis, err := net.Listen("tcp", "localhost:6000")
	if err != nil {
		log.Fatalf("failed to initializa TCP listen: %v", err)
	}
	defer lis.Close()

	server := grpc.NewServer()
	roleServer := &Server{
		userRoles: map[int32][]*pb.Role{
			1: {normal},
			2: {normal, editor},
			3: {normal},
			4: {normal, editor, admin},
			5: {normal, editor, admin, superUser},
		},
		roles: []*pb.Role{normal, editor, admin, superUser},
	}
	pb.RegisterRolesServer(server, roleServer)

	//etcd service configuration
	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}

	//new etcd clients with new configuration
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	//creates new api key
	kapi := client.NewKeysAPI(c)
	// set "/foo" key with "bar" value
	log.Print("Setting '/role' key with 'localhost:6000' value")
	resp, err := kapi.Set(context.Background(), "/role", "localhost:6000", nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Printf("Set is done. Metadata is %q\n", resp)
	}

	server.Serve(lis)

}
