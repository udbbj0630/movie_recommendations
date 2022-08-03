package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	//"test.com/v/movie"

	"sort"

	movie "udbbj0630.github.com/models"
)

type TfxRequest struct {
	SignatureName string
	Instances     [][2]float32
	Inputs        [][]float32
}

type TfxRequestName struct {
	Name []string
}

var testRankMovies = map[string]movie.MetaData{
	"movie1": movie.MetaData{
		Name:        "stimulation 1995",
		VoteAverage: 1,
		VoteCount:   1,
	},
	"movie2": movie.MetaData{
		Name:        "godfather",
		VoteAverage: 2,
		VoteCount:   2,
	},
	"movie3": movie.MetaData{
		Name:        "the American Past",
		VoteAverage: 3,
		VoteCount:   3,
	},
	"movie4": movie.MetaData{
		Name:        "the cinema of heaven",
		VoteAverage: 4,
		VoteCount:   4,
	},
	"movie5": movie.MetaData{
		Name:        "city without a master",
		VoteAverage: 5,
		VoteCount:   5,
	},
}

func mockTfxFailure(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(500)
	fmt.Fprintf(w, "{\"error\": \"Test failure\"}")
}

func mockElasticFailure(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(500)
	fmt.Fprintf(w, "{\"error\": \"Test failure\"}")
}

func mockRankTfxCall(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("-%s\n-", body)
	var tfxRequest []string
	err := json.Unmarshal(body, &tfxRequest)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{\"error\": \"Parsing body failed due to "+err.Error()+"\"")
		return
	}
	var tfxResponse []string
	sort.Strings(tfxRequest)
	for _, name := range tfxRequest {
		fmt.Printf("%s\n", name)
		tfxResponse = append(tfxResponse, name)
	}
	var tfxResponseBytes []byte
	tfxResponseBytes, err = json.Marshal(tfxResponse)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{\"error\": \"Fail to serialize tfx response.\"")
		return
	}
	w.WriteHeader(200)
	w.Write(tfxResponseBytes)
}

func TestMovieRanking(t *testing.T) {
	// Mock TFX server
	tfxServer := httptest.NewServer(http.HandlerFunc(mockRankTfxCall))
	defer tfxServer.Close()

	// Prepare testing recommender.
	rec := recommender{
		Movies: make(map[string]movie.MetaData),
		TfxURL: tfxServer.URL,
	}
	rec.Movies = testRankMovies
	rankingServer := httptest.NewServer(rec.handler())
	defer rankingServer.Close()

	// Make testing call.
	res, err := http.Get(rankingServer.URL + "/movies")
	if err != nil {
		t.Fatal("HTTP call should succeed.")
	}
	if res.StatusCode != 200 {
		t.Fatal("Reqeust failed with status " + res.Status)
	}

	actual, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("Request shouldn't fail")
	}
	fmt.Printf("%s\n", actual)
	expected := []byte("[{\"Name\":\"city without a master\",\"Score\":5},{\"Name\":\"godfather\",\"Score\":4}]")
	for idx, byteValue := range actual {
		if idx == 3 {
			break
		}
		if byteValue != expected[idx] {
			t.Fatal(fmt.Sprintf("Expected response is:\n %s\nBut gotten:\n%s", expected, actual))
		}
	}
}

func TestMovieRankingFailure(t *testing.T) {
	// Prepare TFX mock

	tfxServer := httptest.NewServer(http.HandlerFunc(mockTfxFailure))
	defer tfxServer.Close()

	// Prepare testing recommender.
	rec := recommender{
		Movies: make(map[string]movie.MetaData),
		TfxURL: tfxServer.URL,
	}
	rec.Movies = testRankMovies
	rankingServer := httptest.NewServer(rec.handler())
	defer rankingServer.Close()

	// Make testing call.
	res, err := http.Get(rankingServer.URL + "/movies")
	if err != nil {
		t.Fatal("HTTP call should succeed.")
	}
	if res.StatusCode != 200 {
		return
	}

	actual, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("Request shouldn't fail")
	}
	fmt.Printf("%s\n", actual)
	expected := []byte("[{\"Name\":\"city without a master\",\"Score\":5},{\"Name\":\"godfather\",\"Score\":4}]")
	for idx, byteValue := range actual {
		if idx == 3 {
			break
		}
		if byteValue != expected[idx] {
			t.Fatal(fmt.Sprintf("Expected response is:\n %s\nBut gotten:\n%s", expected, actual))
		}
	}
}
