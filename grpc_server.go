package main

import (
	"context"
	"net"

	pb "go-crud-app/proto"

	"google.golang.org/grpc"
)

type UserGrpcServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserGrpcServer) GetUserById(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	var user User
	DB.First(&user, "id = ?", req.Id)

	return &pb.UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Age:   int32(user.Age),
	}, nil
}

func StartGrpcServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &UserGrpcServer{})

	go grpcServer.Serve(lis)
}
