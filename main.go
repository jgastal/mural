package main

import (
	"html/template"
	"net/http"
	"os"
)

func hello(resp http.ResponseWriter, req *http.Request) {
	content, err := template.New("").Parse("<html><head><title>{{.title}}</title></head><body><ul>{{range .envs}}<li>{{.}}</li>{{end}}</ul></body></html>")
	if err != nil {
		panic(err)
	}

	ctx := map[string]interface{} {
		"title": "How about them apples?!",
		"envs": os.Environ(),
	}

	content.Execute(resp, ctx)
}

func main() {
	http.HandleFunc("/", hello)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
