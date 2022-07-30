package main

import (
	"context"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
	pb "users/grpc-users"
	"users/pkg/db/clickhouse"
	"users/pkg/db/postgres"
	"users/pkg/kafka"
)

const (
	port = ":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (server *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	insertedUser, err := postgres.InsertUser(in.Name, in.Age)
	if err != nil {
		return nil, err
	}

	key, val, err := kafka.Kafka(ctx, insertedUser)
	if err != nil {
		return nil, err
	}

	click, err := clickhouse.Clickhouse()
	if err != nil {
		return nil, err
	}
	_ = click.Create(&clickhouse.Logs{
		Name:      key,
		Age:       val,
		Timestamp: time.Now(),
	})

	return insertedUser, nil
}

func (server *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	users_list, err := postgres.GetUsers()
	if err != nil {
		return nil, err
	}
	return users_list, nil
}

func (server *UserManagementServer) DeleteUser(ctx context.Context, in *pb.DeleteUserParams) (*pb.DeletedUser, error) {

	delUser, err := postgres.DeleteUser(in.Id)
	if err != nil {
		return nil, err
	}
	return delUser, nil

}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}
	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func main() {
	conn, err := postgres.InitDB()
	if err != nil {
		return
	}

	err = postgres.MigrateDB()
	if err != nil {
		log.Printf("Error migrating: %v", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
