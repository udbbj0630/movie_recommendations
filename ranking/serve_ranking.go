package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	resp := make(map[string][]string)
	resp["movieList"] = []string{
		"Nope", 
		"Thor: Love and Thunder",
		"Minions: The Rise of Gru",
		"where the Crawdads Sing",
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error when JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}
