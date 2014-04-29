package main

import (
	"net/http"
	"os"
)

func hello(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("<html><body><h1>Hello world!</h1></body></html>"))
}

func main() {
	http.HandleFunc("/", hello)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
