package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", pingPongHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func pingPongHandler(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Message string `json:"message"`
	}{
		Message: "pong",
	}

	json.NewEncoder(w).Encode(&response)
}
