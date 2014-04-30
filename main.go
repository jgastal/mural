package main

import (
	"database/sql"
	"github.com/lib/pq"
	"html/template"
	"net/http"
	"os"
)

type message struct {
	Name string
	Message string
}

func home(resp http.ResponseWriter, req *http.Request) {
	conn_params, err := pq.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	conn, err := sql.Open("postgres", conn_params)
	if err != nil {
		panic(err)
	}

	if req.Method == "POST" {
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
	}

	home, err := template.ParseFiles("home.html")
	if err != nil {
		panic(err)
	}

	var messages [10]message
	msg_count := 0
	rows, err := conn.Query("SELECT name, message FROM posts ORDER BY id DESC LIMIT 10");
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var name, msg string
		err = rows.Scan(&name, &msg)
		if err != nil {
			//TODO log error
			continue
		}
		messages[msg_count] = message{name, msg}
		msg_count += 1
	}
	conn.Close()

	ctx := map[string]interface{}{
		"title": "How do you like them apples?!",
		"messages": messages[:msg_count],
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
