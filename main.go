package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

const stubLength = 6
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func initRandString() {
	rand.Seed(time.Now().UnixNano())
}

func randString(length int) string {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(bytes)
}

type Entry struct {
	Stub string `json:"stub"`
	Url  string `json:"url"`
}

var entries map[string]string

func setupHeaders(writer *http.ResponseWriter) {
	(*writer).Header().Set("Access-Control-Allow-Origin", "*")
	(*writer).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*writer).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*writer).Header().Add("Content-Type", "application/json")
}

func newShortUrl(writer http.ResponseWriter, req *http.Request) {
	setupHeaders(&writer)

	var entry Entry
	err := json.NewDecoder(req.Body).Decode(&entry)
	if err != nil {
		log.Fatal(err)
	}

	if entry.Stub == "" || entries[entry.Stub] != "" {
		entry.Stub = randString(stubLength)
	}

	entries[entry.Stub] = entry.Url

	json.NewEncoder(writer).Encode(entry)
}

func getShortUrl(writer http.ResponseWriter, req *http.Request) {
	setupHeaders(&writer)

	var entry Entry
	entry.Stub = mux.Vars(req)["stub"]

	entry.Url = entries[entry.Stub]

	json.NewEncoder(writer).Encode(entry)
}

func main() {
	entries = make(map[string]string)

	log.Println("Reading PORT from Environment")
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT must be provided")
	}

	log.Println("Initializing Random String Generator")
	initRandString()

	log.Println("Initializing Router")
	router := mux.NewRouter()

	log.Println("Setting up Routes")
	router.HandleFunc("/api/shorten", newShortUrl).Methods("POST")
	log.Println("[+] POST: /api/shorten")
	router.HandleFunc("/api/shorten/{stub}", getShortUrl).Methods("GET")
	log.Println("[+] GET: /api/shorten/{stub}")

	log.Printf("Starting server on %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
