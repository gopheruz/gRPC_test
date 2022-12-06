package main

import (
	"context"
	"fmt"
	service "gRPCpostges/genproto/postgres_service"
	"log"
	"net"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	service.UnimplementedUserSServer
	db *sqlx.DB
}

var (
	Host     string = "localhost"
	Port     string = "5432"
	User     string = "postgres"
	Password string = "postgrespw"
	Database string = "grpc_test"
)

func main() {
	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		Host,
		Port,
		User,
		Password,
		Database,
	)
	psqlConn, err := sqlx.Connect("postgres", psqlUrl)

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	fmt.Println(psqlUrl)
	fmt.Println("Connected Succesfully!")

	listener, err := net.Listen("tcp", ":8002")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	service.RegisterUserSServer(srv, &Server{db: psqlConn})
	reflection.Register(srv)

	fmt.Println("Server Started")

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}

func (user *Server) Create(ctx context.Context, userReq *service.UserReq) (*service.UserRes, error) {
	var result service.UserRes
	query := `
		INSERT INTO users (
			first_name,
			last_name,
			age
		) VALUES  ($1, $2, $3)
		RETURNING id, first_name, last_name, age
	`
	row := user.db.QueryRow(query,
		userReq.FirstName,
		userReq.LastName,
		userReq.Age,
	)
	err := row.Scan(
		&result.Id,
		&result.FirstName,
		&result.LastName,
		&result.Age,
	)
	if err != nil {
		return nil, err
	}
	return &result, err
}

func (user *Server) Update(ctx context.Context, userReq *service.UserReq) (*service.UserRes, error) {
	var result service.UserRes
	query := `
		UPDATE users SET  
			first_name = $1,
			last_name = $2,
			age = $3
		WHERE id = $4
		RETURNING id, first_name, last_name, age
	`
	row := user.db.QueryRow(query,
		userReq.FirstName,
		userReq.LastName,
		userReq.Age,
		userReq.Id,
	)
	err := row.Scan(
		&result.Id,
		&result.FirstName,
		&result.LastName,
		&result.Age,
	)
	if err != nil {
		return nil, err
	}
	return &result, err
}

func (user *Server) Get(ctx context.Context, userReq *service.IdMsg) (*service.UserRes, error) {
	var result service.UserRes
	query := "SELECT first_name, last_name, age FROM users where id = $1"
	row := user.db.QueryRow(query, userReq.Id)
	err := row.Scan(
		&result.FirstName,
		&result.LastName,
		&result.Age,
	)
	result.Id = userReq.Id
	if err != nil {
		return nil, err
	}
	return &result, err
}

func (user *Server) Delete(ctx context.Context, userReq *service.IdMsg) (*service.Empty, error) {
	query := "DELETE FROM users where id = $1"
	_, err := user.db.Exec(query, userReq.Id)
	if err != nil {
		return nil, err
	}
	return &service.Empty{}, nil
}

func (user *Server) GetAll(ctx context.Context, params *service.GetAlluserParams) (*service.GetAllUserResponse, error) {
	offset := (params.Page - 1) * params.Limit
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)
	filter := ""
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter = fmt.Sprintf(" WHERE first_name ILIKE '%s' OR last_name ILIKE '%s'", str, str)
	}
	query := `
		SELECT 
			id,
			first_name,
			last_name,
			age
		FROM users 
	` + filter + limit
	rows, err := user.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res service.GetAllUserResponse
	res.Users = make([]*service.UserRes, 0)
	for rows.Next() {
		var u service.UserRes
		err := rows.Scan(
			&u.Id,
			&u.FirstName,
			&u.LastName,
			&u.Age,
		)
		if err != nil {
			return nil, err
		}
		res.Users = append(res.Users, &u)
	}
	queryCount := "SELECT COUNT(1) FROM users" + filter
	err = user.db.QueryRow(queryCount).Scan(&res.Count)
	if err != nil {
		return nil, err
	}
	return &res, nil

}
