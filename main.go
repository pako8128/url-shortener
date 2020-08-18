package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Entry struct {
	Stub string `json:"stub"`
	Url  string `json:"url"`
}

var entries map[string]string

func newShortUrl(writer http.ResponseWriter, req *http.Request) {
	var entry Entry
	err := json.NewDecoder(req.Body).Decode(&entry)
	if err != nil {
		log.Fatal(err)
	}

	entries[entry.Stub] = entry.Url

	writer.Header().Add("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(entry)
}

func getShortUrl(writer http.ResponseWriter, req *http.Request) {
	var entry Entry
	entry.Stub = mux.Vars(req)["stub"]

	entry.Url = entries[entry.Stub]

	writer.Header().Add("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(entry)
}

func main() {
	entries = make(map[string]string);

	log.Println("Initializing Router")
	router := mux.NewRouter()

	log.Println("Setting up Routes")
	router.HandleFunc("/api/shorten", newShortUrl).Methods("POST")
	log.Println("[+] POST: /api/shorten")
	router.HandleFunc("/api/shorten/{stub}", getShortUrl).Methods("GET")
	log.Println("[+] GET: /api/shorten/{stub}")

	port := ":8000"

	log.Printf("Starting server on %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
