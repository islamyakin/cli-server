package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "cli-server/auth"
)

const (
	port       = ":50051"
	validToken = "secret-token"
	validUser  = "admin"
	validPass  = "password"
)

type server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	if in.GetUsername() == validUser && in.GetPassword() == validPass {
		return &pb.LoginResponse{Token: validToken, Message: "Login berhasil, Halo " + validUser}, nil
	}
	return nil, fmt.Errorf("invalid credentials")
}

func (s *server) SendMessage(ctx context.Context, in *pb.MessageRequest) (*pb.MessageResponse, error) {
	if in.GetToken() != validToken {
		return nil, fmt.Errorf("unauthorized")
	}
	log.Println("Received message from client:", in.GetMessage())
	return &pb.MessageResponse{Response: "Message received: " + in.GetMessage()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{})
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
