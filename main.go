package main

import (
	"fmt"
	"net/http"
	"os"
)

func hello(resp http.ResponseWriter, req *http.Request) {
	fmt.Println(resp, "<html><head><title>How about them apples?!</title></head><body><h1>")
	fmt.Println(resp, "Hello world!")
	fmt.Println(resp, "</h1></body></html>")
}

func main() {
	http.HandleFunc("/", hello)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
