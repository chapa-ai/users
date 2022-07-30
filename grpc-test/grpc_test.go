package test

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"testing"
	pb "users/grpc-users"
)

func TestCreateNewUser(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewUserManagementClient(conn)
	_, err = client.CreateNewUser(ctx, &pb.NewUser{Name: "Krylov", Age: "17"})
	if err != nil {
		t.Fatalf("CreateNewUser failed: %v", err)
	}
}

func TestDeleteUsers(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewUserManagementClient(conn)
	_, err = client.DeleteUser(ctx, &pb.DeleteUserParams{Id: 10})
	if err != nil {
		t.Fatalf("DeleteUsers failed: %v", err)
	}
}

func TestGetUsers(t *testing.T) {
	ctx := context.Background()

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewUserManagementClient(conn)
	_, err = client.GetUsers(ctx, &pb.GetUsersParams{})
	if err != nil {
		t.Fatalf("GetUsers failed: %v", err)
	}
}
