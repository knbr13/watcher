package main

import (
	"fmt"
	"log"
	"net/http"
)

// this is a demo server to try the watcher tool on it

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world")
	})
	log.Println("starting...")
	http.ListenAndServe(":8080", nil)
}
