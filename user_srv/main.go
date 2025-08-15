package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"mxshop/user_srv/handler"
	"mxshop/user_srv/proto"
	"net"
)

func main() {
	IP := flag.String("ip", "127.0.0.1", "ip address")
	PORT := flag.Int("port", 50051, "port number")
	flag.Parse()
	path := fmt.Sprintf("%s:%d", *IP, *PORT)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", path)
	if err != nil {
		panic(err)
	}
	err = server.Serve(lis)
	if err != nil {
		panic(err)
	}
}
