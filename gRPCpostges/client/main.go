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


	data, err := client.Get(context.Background(), &service.IdMsg{Id: 3})
	if err != nil {
		panic(err)
	}
	fmt.Println(data)

	
	data, err = client.Update(context.Background(), &service.UserReq{
		FirstName: faker.FirstName(),
		LastName:  faker.LastName(),
		Age:       3,
		Id:        3,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
	many, err := client.GetAll(context.Background(), &service.GetAlluserParams{
		SortBy: "id",
		Page:   1,
		Limit:  10,
	})
	if err != nil {
		panic(err)
	}
	fmt.Print("\n\n###################GETALL################\n")
	for _, val := range many.Users {
		fmt.Println(val)
	}
	fmt.Printf("count %d, \n\n", many.Count)


	_, err= client.Delete(context.Background(), &service.IdMsg{Id: 5})
	if err != nil{
		panic(err)
	}

}
