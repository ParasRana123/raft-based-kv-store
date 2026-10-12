package main

import (
	"log"
	"net/http"

	"raft-kv-store/api"
	"raft-kv-store/kv"
)

func main() {
	store := kv.NewStore()
	server := api.NewServer(store)

	addr := ":8080"

	log.Printf("KV server starting on %s", addr)

	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatal(err)
	}
}
