package main

import (
	"html/template"
	"labix.org/v2/mgo"
	"net/http"
	"os"
	"time"
)

type message struct {
	Name string
	Message string
	Time time.Time
}

func home(resp http.ResponseWriter, req *http.Request) {
	session, err := mgo.Dial(os.Getenv("MONGOHQ_URL"))
	if err != nil {
		panic(err)
	}
	db := session.DB("")
	collection := db.C("posts")

	if req.Method == "POST" {
		name := req.FormValue("name")
		if name == "" {
			name = "Anonymous coward"
		}
		msg := req.FormValue("message")
		if msg == "" {
			msg = "Gosh, nothing to say?!"
		}
		err = collection.Insert(message{name, msg, time.Now()})
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
	iter := collection.Find(nil).Sort("-time").Batch(10).Iter()
	for iter.Next(&messages[msg_count]) && msg_count < 10 {
		msg_count += 1
	}
	if err := iter.Close(); err != nil {
		panic(err)
	}
	session.Close()

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
