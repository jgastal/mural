package main

import (
	"fmt"
	"net/http"
	"os"
)

func hello(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("<html><body><h1>")
	resp.Write([]byte("Hello world!"))
	fmt.Println("</h1></body></html>")
}

func main() {
	http.HandleFunc("/", hello)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
