package main

import (
	"context"
	"fmt"
	service "gRPCpostges/genproto/postgres_service"

	"github.com/bxcodec/faker/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := service.NewUserSClient(conn)
	user, err := client.Create(context.Background(), &service.UserReq{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Age:       18,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
	user, err = client.Update(context.Background(), &service.UserReq{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Age:       18,
		Id:        5,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
	_, err = client.Delete(context.Background(), &service.IdMsg{Id: 9})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("user deleted succesfully")

	data, err := client.GetAll(context.Background(), &service.GetAlluserParams{
		Limit:  10,
		Page:   1,
		Search: "a",
		SortBy: "id",
	})
	if err != nil {
		panic(err)
	}
	for _, val := range data.Users {
		fmt.Println(val)
	}
	fmt.Println(data.Count)
}
