package main

import (
	"database/sql"
	"github.com/lib/pq"
	"html/template"
	"net/http"
	"os"
)

func home(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		conn_params, err := pq.ParseURL(os.Getenv("DATABASE_URL"))
		if err != nil {
			panic(err)
		}
		conn, err := sql.Open("postgres", conn_params)
		if err != nil {
			panic(err)
		}

		name := req.FormValue("name")
		if name == "" {
			name = "Anonymous coward"
		}
		message := req.FormValue("message")
		if message == "" {
			message = "Gosh, nothing to say?!"
		}
		_, err = conn.Exec("INSERT INTO posts (name, message) VALUES ($1, $2)", name, message)
		if err != nil {
			panic(err)
		}

		conn.Close()
	}

	home, err := template.ParseFiles("home.html")
	if err != nil {
		panic(err)
	}

	ctx := map[string]interface{}{
		"title": "How about them apples?!",
		"envs":  os.Environ(),
	}

	home.Execute(resp, ctx)
	
}

func main() {
	http.HandleFunc("/", home)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
