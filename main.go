package main

import (
	"container/list"
	"github.com/gorilla/websocket"
	"html/template"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"os"
	"time"
)

var upgrader = websocket.Upgrader{}

var clients = list.New()

type message struct {
	Name    string
	Message string
	Time    time.Time
}

func newClient(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	el := clients.PushBack(ws)
	log.Printf("Got a new client %p, now have: %d clients", ws, clients.Len())

	// Monitor the connection and remove from clients list when it closes
	for {
		if err = ws.ReadJSON(nil); err != nil {
			clients.Remove(el)
		}
	}
}

func notifyUsers(msg message) {
	for e := clients.Front(); e != nil; e = e.Next() {
		ws := e.Value.(*websocket.Conn)
		log.Printf("Notifying user %p of new message", ws)
		go websocket.WriteJSON(ws, msg)
	}
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
		m := message{name, msg, time.Now()}
		err = collection.Insert(m)
		if err != nil {
			panic(err)
		}
		session.Close()
		go notifyUsers(m)
		return
	}

	home, err := template.ParseFiles("home.html")
	if err != nil {
		panic(err)
	}

	var messages [10]message
	msg_count := 0
	iter := collection.Find(nil).Sort("-time").Batch(10).Iter()
	for msg_count < 10 && iter.Next(&messages[msg_count]) {
		msg_count += 1
	}
	if err := iter.Close(); err != nil {
		panic(err)
	}
	session.Close()

	ctx := map[string]interface{}{
		"title":    "How do you like them apples?!",
		"messages": messages[:msg_count],
	}

	home.Execute(resp, ctx)
}

func main() {
	http.HandleFunc("/websocket/", newClient)

	http.HandleFunc("/msg.mst", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "msg.mst")
	})
	http.HandleFunc("/", home)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
