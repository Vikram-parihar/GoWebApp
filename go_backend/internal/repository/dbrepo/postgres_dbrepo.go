package dbrepo

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PostgresDBRepo struct {
	DB *sql.DB
}
const dbTimeout= time.Second*3

func (m *PostgresDBRepo)Connection() *sql.DB{
	return m.DB
}

/*
The context package in Go provides a way to carry request-scoped values across API boundaries and between processes. It is widely used in Go programming for managing and canceling requests, timeouts, and deadlines in a clean and efficient way.
The context package provides the Context type, which represents the context of a request or a task. It contains a Context value that can be used to store key-value pairs of request-scoped data, such as user authentication information, request ID, or other metadata.
Some of the main features of the context package are:
Contexts can be canceled, which is useful for canceling long-running tasks or for implementing timeouts.
A context can have parent contexts, allowing a request to be broken down into smaller, nested contexts.
Values can be retrieved from the context using a key and value pair, providing a convenient way to pass data between different layers of an application.
Contexts can be used to propagate metadata across API boundaries, such as request headers or authentication tokens.
Contexts are safe for concurrent use, making them suitable for use in concurrent and distributed systems.
The context package is a powerful tool for writing clean, efficient, and scalable Go programs. It is widely used in the Go standard library and in many third-party libraries and frameworks.
*/
func (m *PostgresDBRepo) AllMovies() ([]*models.Movie,error){

	

	ctx,cancel:=context.WithTimeout(context.Background(),dbTimeout)
	defer cancel()

	query:=`
		select 
		  id, title, release_date, runtime,
		  mpaa_rating, description, coalesce(image, ''),
		  created_at, updated_at 
		from 
			movies
		order by 
			title
		`
		rows,err:= m.DB.QueryContext(ctx, query)
		if err!=nil{
			return nil, err
		}
		defer rows.Close()


	var movies []*models.Movie
	for rows.Next(){
		var movie models.Movie
		err:=rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.RunTime,
			&movie.MPAARating,
			&movie.Description,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err !=nil{
			return nil,err
		}
		movies=append(movies, &movie)
	}
	return movies, nil
}

func (m *PostgresDBRepo) GetUserByEmail(email string)(*models.User,error){
	ctx,cancel:=context.WithTimeout(context.Background(),dbTimeout)
	defer cancel()

	query:=`select id,email,first_name, last_name,password,
	        created_at,updated_at from users where email= $1`
	var User models.User
	row:=m.DB.QueryRowContext(ctx, query,email)
	err:=row.Scan(
		&User.ID,
		&User.Email,
		&User.FirstName,
		&User.LastName,
		&User.Password,
		&User.CreatedAt,
		&User.UpdatedAt,
	)
	if err !=nil{
		fmt.Println(err)
		return nil,err
	}
	fmt.Println(User)
	
	return &User,nil

}

func (m *PostgresDBRepo) GetUserByID(id int)(*models.User,error){
	ctx,cancel:=context.WithTimeout(context.Background(),dbTimeout)
	defer cancel()

	query:=`select id,email,first_name, last_name,password,
	        created_at,updated_at from users where id= $1`
	var User models.User
	row:=m.DB.QueryRowContext(ctx, query,id)
	err:=row.Scan(
		&User.ID,
		&User.Email,
		&User.FirstName,
		&User.LastName,
		&User.Password,
		&User.CreatedAt,
		&User.UpdatedAt,
	)
	if err !=nil{
		fmt.Println(err)
		return nil,err
	}
	fmt.Println(User)
	
	return &User,nil

}