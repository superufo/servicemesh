package main

import (
	"log"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/server"
	pb "github.com/go-chassis/ygx/sidecar/example/grpc/helloworld"

	_ "github.com/go-chassis/go-chassis/server/grpc"
	"github.com/go-mesh/openlogging"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// Server if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/grpc/server/
// Server is used to implement helloworld.GreeterServer.
type Server struct{}

// SayHello implements helloworld.GreeterServer
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	log.Println(md["x-user"])
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	chassis.RegisterSchema("grpc", &Server{}, server.WithRPCServiceDesc(&pb.Greeter_serviceDesc))
	if err := chassis.Init(); err != nil {
		openlogging.Error("Init failed.")
		return
	}
	chassis.Run()
}
