package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Print("Starting server on :3000")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
