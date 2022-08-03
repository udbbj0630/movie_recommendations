package movie

import (
	"encoding/csv"
	"strings"
	"testing"
)

var header = strings.Join([]string{
	"adult", "belongs_to_collection", "budget", "genres", 
	"homepage", "id", "imdb_id", "original_language", 
	"original_title", "overview", "popularity", "poster_path",
	"production_companies", "production_countries", "release_date", "revenue",
	"runtime", "spoken_languages", "status", "tagline",
	"title", "video", "vote_average", "vote_count"}, ",")
var goodLine = strings.Join([]string{
	"true", "", "10000", "",
	"", "", "", "en",
	"title", "overview", "", "",
	"", "", "", "",
	"", "", "", "", 
	"", "", "5.1", "1000"}, ",")
var badLine = strings.Join([]string{
	"true", "", "10000", "",
	"", "", "", // Wrong column number by missing language "en",
	"title", "overview", "", "",
	"", "", "", "",
	"", "", "", "", 
	"", "", "5.1", "1000"}, ",")
var badVoteCount = strings.Join([]string{
	"true", "", "10000", "",
	"", "", "", // Wrong column number by missing language "en",
	"title", "overview", "", "",
	"", "", "", "",
	"", "", "", "", 
	"", "", "5.1", "a1000"}, ",")
var badVoteAverage = strings.Join([]string{
	"true", "", "10000", "",
	"", "", "", // Wrong column number by missing language "en",
	"title", "overview", "", "",
	"", "", "", "",
	"", "", "", "", 
	"", "", "a5.1", "1000"}, ",")

func TestBuildMoviesMetadataSuccess(t *testing.T) {
	r := csv.NewReader(strings.NewReader(strings.Join([]string{header, goodLine}, "\n")))
	built := make(map[string]MetaData)
	err := BuildMoviesMetadata(r, &built)
	if err != nil {
		t.Fatal("Parsing valid csv entry failed.")
	}
	if (len(built)!= 1) {
		t.Fatal("Fails to parse valid line")
	}
}

func TestBuildMoviesMetadataSkipsBadLine(t *testing.T) {
	r := csv.NewReader(strings.NewReader(strings.Join([]string{header, badLine}, "\n")))
	built := make(map[string]MetaData)
	err := BuildMoviesMetadata(r, &built)
	if (err != nil) {
		t.Fatal("Fail to skip invalid csv entry.")
	}
	if (len(built) != 0) {
		t.Fatal("Invalid csv entry isn't skpped.")
	}
}

func TestBuildMoviesMetadataHeaderMustExits(t *testing.T) {
	r := csv.NewReader(strings.NewReader(""))
	built := make(map[string]MetaData)
	err := BuildMoviesMetadata(r, &built)
	if (err == nil) {
		t.Fatal("Bad header file shouldn't be included.")
	}
	if (len(built) != 0) {
		t.Fatal("Invalid csv entry isn't skpped.")
	}
}

func TestBuildMoviesMetadataSkipsBadVoteCount(t *testing.T) {
	r := csv.NewReader(strings.NewReader(strings.Join([]string{header, badVoteCount}, "\n")))
	built := make(map[string]MetaData)
	err := BuildMoviesMetadata(r, &built)
	if (err != nil) {
		t.Fatal("Fail to skip bad entry")
	}
	if (len(built) != 0) {
		t.Fatal("Fail to skip bad vote count")
	}
}

func TestBuildMoviesMetadataSkipsBadVoteAverage(t *testing.T) {
	r := csv.NewReader(strings.NewReader(strings.Join([]string{header, badVoteAverage}, "\n")))
	built := make(map[string]MetaData)
	err := BuildMoviesMetadata(r, &built)
	if (err != nil) {
		t.Fatal("Fail to skip bad entry")
	}
	if (len(built) != 0) {
		t.Fatal("Fail to skip bad vote average")
	}
}