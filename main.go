package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Initializing Router")
	router := mux.NewRouter()

	port := ":8000"

	log.Printf("Starting server on %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
