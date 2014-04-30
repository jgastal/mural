package main

import (
	"database/sql"
	"github.com/lib/pq"
	"html/template"
	"net/http"
	"os"
)

func hello(resp http.ResponseWriter, req *http.Request) {
	conn_params, err := pq.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	conn, err := sql.Open("postgres", conn_params)
	if err != nil {
		panic(err)
	}

	home, err := template.ParseFiles("home.html")
	if err != nil {
		panic(err)
	}

	ctx := map[string]interface{}{
		"title": "How about them apples?!",
		"envs":  os.Environ(),
	}
	if conn.Ping() != nil {
		ctx["ps"] = "I connected to the database! YAY!"
	} else {
		ctx["ps"] = "I can't connect to the database! WTF?!"
	}

	home.Execute(resp, ctx)
}

func main() {
	http.HandleFunc("/", hello)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
