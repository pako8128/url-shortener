package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

func newShortUrl(writer http.ResponseWriter, req *http.Request) {
	var entry Entry
	err := json.NewDecoder(req.Body).Decode(&entry)
	if err != nil {
		log.Fatal(err)
	}

	if entry.Stub == "" || entries[entry.Stub] != "" {
		entry.Stub = randString(stubLength)
	}

	entries[entry.Stub] = entry.Url

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(entry)
}

func getShortUrl(writer http.ResponseWriter, req *http.Request) {
	stub := mux.Vars(req)["stub"]
	url := entries[stub]

	if url == "" {
		http.Redirect(writer, req, "https://url-shortener-kohl.vercel.app/", http.StatusSeeOther)
	} else {
		http.Redirect(writer, req, url, http.StatusSeeOther)
	}
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
	router.HandleFunc("/{stub}", getShortUrl).Methods("GET")
	log.Println("[+] GET: /{stub}")

	log.Println("Wrapping Router with CORS Handler")
	wrapper := cors.New(cors.Options{
		AllowedOrigins: []string{"https://url-shortener-kohl.vercel.app", "https://url-shortener.pako8128.vercel.app"},
		AllowedHeaders: []string{"Content-Type", "Accept-Encoding", "Accept-Language", "DNT"},
		Debug:          true,
	})
	handler := wrapper.Handler(router)

	log.Printf("Starting server on %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
