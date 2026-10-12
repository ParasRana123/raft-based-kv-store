package main

import (
	"log"
	"net/http"
	"os"

	"raft-kv-store/api"
	"raft-kv-store/kv"
)

func main() {
	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		dataFile = "data.json"
	}

	store, err := kv.NewPersistentStore(dataFile)
	if err != nil {
		log.Fatalf("Failed to initialize persistent store (%s): %v", dataFile, err)
	}

	server := api.NewServer(store)

	addr := ":8080"
	log.Printf("KV server starting on %s with persistence file: %s", addr, dataFile)

	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatal(err)
	}
}
