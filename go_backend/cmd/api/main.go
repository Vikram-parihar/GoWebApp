package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const port = 8080

type application struct {
	DSN    string
	Domain string
	DB     repository.DatabaseRepo
	auth Auth
	JWTSecret string
	JWTIssuer string
	JWTAudience string 
	CookieDomain string
}

func main() {
	var app application
	log.Println("Application is starting at:", port)
	// helloHandler := func(w http.ResponseWriter, req *http.Request) {
	// 	fmt.Fprint(w,"Hello, world!\n")
	// }
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.JWTSecret ,"jwt-secret","verysecret","Signing secret")
	flag.StringVar(&app.JWTIssuer ,"jwt-issuer","example.com","Signing issuer")
	flag.StringVar(&app.JWTAudience ,"jwt-audience","example.com","Signing audience")
	flag.StringVar(&app.CookieDomain ,"cookie-domain","localhost","cookie domain")
	flag.StringVar(&app.Domain ,"domain","example.com","domain")

	flag.Parse()

	//connect to a database
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()
	// // http.HandleFunc("/",helloHandler)
	app.auth=Auth{
		Issuer: app.JWTIssuer,
		Audience: app.JWTAudience,
		Secret: app.JWTSecret,
		TokenExpiry: time.Minute*15,
		RefreshExpiry: time.Hour*24,
		CookiePath: "/",
		CookieName: "refresh_token",
		CookieDomain: app.CookieDomain,
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
