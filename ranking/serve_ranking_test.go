package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

var buff []byte

type TestResponseWriter struct {
	statusCode int;
	header map[string][]string; 
}

func (t TestResponseWriter) Header() http.Header {
	header := http.Header{}
	return header
}

func (t TestResponseWriter) Write(b []byte) (int, error) {
	buff = b
	return len(b), nil
}

func (t TestResponseWriter) WriteHeader(statusCode int) {
	t.statusCode = statusCode
}

func TestHandler(t *testing.T) {
	writer := TestResponseWriter{}
	req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
	handler(writer, req)

	gotJson := buff[:]
	got := make(map[string][]string)
	if err := json.Unmarshal(gotJson, &got); err != nil {
		t.Fatalf("unmarshal resp error: %s", err)
	}
	movieList := got["movieList"]
	wantMovieList := []string{
		"Nope", 
		"Thor: Love and Thunder",
		"Minions: The Rise of Gru",
		"where the Crawdads Sing",
	}
	for i, movie := range wantMovieList {
		if movie != movieList[i] {
			t.Fatalf("Expected %q, but got %q", movie, movieList[i])
		}
	}
}
