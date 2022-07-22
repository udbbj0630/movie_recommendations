package main

import (
	"testing"
	"net/http"
)
var buff []byte

type TestResposeWriter struct{
} 
func (t TestResposeWriter) Header() http.Header {
	return nil
}
func (t TestResposeWriter) Write(b []byte) (int, error) {
	buff = b
	return 0,nil
}
func (t TestResposeWriter) WriteHeader(statusCode int) {
}


func TestHandler(t *testing.T) {
	writer := TestResposeWriter{}
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	handler(writer, req)

	want := "Movie List Nope,Thor: Love and Thunder,Minions: The Rise of Gru,Where the Crawdads Sing"
	got := string(buff[:])
	if got != want 	{ 
		t.Fatalf(" expected %q, but got %q", want, got)
	}
}