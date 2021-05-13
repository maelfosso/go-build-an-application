package main

import (
	"log"
	"net/http"
)

type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 123
}

func main() {
	server := &PlayerServer{&InMemoryPlayerStore{}}
	// handler := http.HandlerFunc(PlayerServer)
	log.Fatal(http.ListenAndServe(":5000", server))
}
