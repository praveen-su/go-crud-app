package main

import (
	"context"
	"net"

	pb "go-crud-app/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

func (s *UserGrpcServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {
	task := Task{
		ID:     uuid.New().String(),
		UserID: req.UserId,
		Title:  req.Title,
		Status: "OPEN",
	}

	DB.Create(&task)

	return &pb.TaskResponse{
		Id:     task.ID,
		UserId: task.UserID,
		Title:  task.Title,
		Status: task.Status,
	}, nil
}

func (s *UserGrpcServer) CreateMultipleTasks(ctx context.Context, req *pb.TaskList) (*pb.TaskList, error) {
	for _, t := range req.Tasks {
		task := Task{
			ID:     uuid.New().String(),
			UserID: t.UserId,
			Title:  t.Title,
			Status: "OPEN",
		}
		DB.Create(&task)
	}
	return req, nil
}

func StartGrpcServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, &UserGrpcServer{})
	reflection.Register(server)
	go server.Serve(lis)
}
