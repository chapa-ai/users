package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc"
	"log"
	"strconv"
	pb "users/grpc-users"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < 100; i++ {
		age := strconv.Itoa(gofakeit.Number(1, 70))
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: gofakeit.Name(), Age: age})
		if err != nil {
			log.Fatalf("could not create user: %v", err)
		}
		logUser := fmt.Sprintf(`User Details: NAME: %s AGE: %s ID: %d`, r.GetName(), r.GetAge(), r.GetId())
		log.Print(logUser)
	}

	params := &pb.GetUsersParams{}
	a, err := c.GetUsers(ctx, params)
	if err != nil {
		log.Fatalf("could not retrieve users: %v", err)
	}
	log.Print("\nUSER LIST: \n")
	fmt.Printf("a.GetUsers(): %v\n", a.GetUsers())

	_, err = c.DeleteUser(ctx, &pb.DeleteUserParams{Id: 1})
}
