package main

import (
	"log"
	"net/http"
	"os"

	"raft-kv-store/api"
	"raft-kv-store/kv"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	store := kv.NewStore()
	server := api.NewServer(store)

	addr := ":" + port
	log.Printf("Starting Raft KV Store HTTP server on %s", addr)
	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
