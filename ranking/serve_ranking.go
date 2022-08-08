package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"

	//"test.com/v/movie"
	"io/ioutil"

	movie "udbbj0630.github.com/movie_recommendation/models"
)

// Elastic request query body
type elasticQueryString struct {
	Query string `json:"query"`
}

type elasticQuery struct {
	QueryString elasticQueryString `json:"query_string"`
}

type elasticQueryBody struct {
	Source []string     `json:"_source"`
	Query  elasticQuery `json:"query"`
}

// ElasticResponse represents elastic search query response
type elasticResponseShards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type elasticTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type elasticSource struct {
	Name string `json:"name"`
}

type elasticHit struct {
	Index  string        `json:"_index"`
	Type   string        `json:"_type"`
	ID     string        `json:"_id"`
	Score  float32       `json:"_score"`
	Source elasticSource `json:"_source"`
}

type elasticResponseHits struct {
	Total    elasticTotal `json:"total"`
	MaxScore float32      `json:"max_score"`
	Hits     []elasticHit `json:"hits"`
}

type elasticResponse struct {
	Took    int                   `json:"took"`
	TimeOut bool                  `json:"time_out"`
	Shards  elasticResponseShards `json:"_shards"`
	Hits    elasticResponseHits   `json:"hits"`
}

// Recommender holds parameters for ranking handler functions.
type recommender struct {
	TfxURL     string
	ElasticURL string
	Movies     map[string]movie.MetaData
}

func (rec *recommender) buildInstances() ([][2]float32, []string) {
	// Prepare ranking requests for movies.
	instances := make([][2]float32, len(rec.Movies))
	requestOrder := make([]string, len(rec.Movies))
	movieIdx := 0
	for movieName, movieRecord := range rec.Movies {
		instances[movieIdx] =
			[2]float32{
				movieRecord.VoteAverage,
				movieRecord.VoteCount}
		requestOrder[movieIdx] = movieName
		movieIdx++
	}
	return instances, requestOrder
}

func (rec *recommender) buildInstancesName() []string {
	// Prepare ranking requests for movies.

	requestOrder := make([]string, len(rec.Movies))
	movieIdx := 0
	for _, movieRecord := range rec.Movies {

		requestOrder[movieIdx] = movieRecord.Name
		movieIdx++
	}
	return requestOrder
}
func (rec *recommender) buildInstancesFromMovieList(requestOrder []string) [][2]float32 {
	instances := make([][2]float32, len(requestOrder))
	for idx, movie := range requestOrder {
		instances[idx] = [2]float32{
			rec.Movies[movie].VoteAverage,
			rec.Movies[movie].VoteCount,
		}
	}
	return instances
}

func (rec *recommender) movieRanking(w http.ResponseWriter, req *http.Request) {
	// Call TFX to get movie ranking scores.
	requestOrder := rec.buildInstancesName()
	// Prepare TFX request.

	requestBody, err := json.Marshal(requestOrder)
	fmt.Println(string(requestBody))
	if err != nil {
		fmt.Fprintf(w, "Error: tfx request failed due to "+err.Error())
		w.WriteHeader(500)
		return
	}
	// Make TFX calls.
	ranked, err := http.Post(rec.TfxURL+"/predict", "Content-Type: application/json; charset=utf-8", bytes.NewBuffer(requestBody))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: TFX request failed due to "+err.Error())
		return
	}
	body, _ := ioutil.ReadAll(ranked.Body)
	if ranked.StatusCode != 200 {
		w.WriteHeader(500)
		fmt.Fprintf(w, "ERROR: TFX call failed for "+string(body))
		return
	}
	// Load TFX response.
	var rankResult []string
	json.Unmarshal(body, &rankResult)
	fmt.Printf("debug: ")
	fmt.Println(len(rankResult))
	rankedMovies := make([]movie.Result, len(rankResult))
	for idx, name := range rankResult {
		rankedMovies[idx].Name = name
		rankedMovies[idx].Score = float32(len(rankResult) - idx)
	}
	// Sore ranked movies.
	sort.Slice(rankedMovies, func(i, j int) bool {
		return rankedMovies[i].Score > rankedMovies[j].Score
	})
	responseContent, err := json.Marshal(rankedMovies)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: Serialize ranked result failed due to "+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseContent)
}

func (rec *recommender) movieRanking0(w http.ResponseWriter, req *http.Request) {
	// Call TFX to get movie ranking scores.
	instances, requestOrder := rec.buildInstances()
	// Prepare TFX request.
	request := make(map[string][][2]float32)
	request["instances"] = instances
	requestBody, err := json.Marshal(request)
	fmt.Println(string(requestBody))
	if err != nil {
		fmt.Fprintf(w, "Error: tfx request failed due to "+err.Error())
		w.WriteHeader(500)
		return
	}
	// Make TFX calls.
	ranked, err := http.Post(rec.TfxURL+"/predict", "Content-Type: application/json; charset=utf-8", bytes.NewBuffer(requestBody))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: TFX request failed due to "+err.Error())
		return
	}
	body, _ := ioutil.ReadAll(ranked.Body)
	if ranked.StatusCode != 200 {
		w.WriteHeader(500)
		fmt.Fprintf(w, "ERROR: TFX call failed for "+string(body))
		return
	}
	// Load TFX response.
	rankResult := make(map[string][][]float32)
	json.Unmarshal(body, &rankResult)
	fmt.Printf("debug: ")
	fmt.Println(len(rankResult["predictions"]))
	rankedMovies := make([]movie.Result, len(rankResult["predictions"]))
	for idx, score := range rankResult["predictions"] {
		rankedMovies[idx].Name = requestOrder[idx]
		rankedMovies[idx].Score = score[0]
	}
	// Sore ranked movies.
	sort.Slice(rankedMovies, func(i, j int) bool {
		return rankedMovies[i].Score < rankedMovies[j].Score
	})
	responseContent, err := json.Marshal(rankedMovies)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: Serialize ranked result failed due to "+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseContent)
}

func (rec *recommender) autocomplete(w http.ResponseWriter, req *http.Request) {
	// Extract query.
	q, ok := req.URL.Query()["q"]
	if !ok || len(q) < 1 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "{\"error\": \"Qeury must contain a param q.\"}")
		return
	}
	query := q[0]
	// Build elastic query request.
	elasticQuery := elasticQueryBody{
		Source: []string{"name"},
		Query: elasticQuery{
			QueryString: elasticQueryString{
				Query: query,
			},
		},
	}
	elasticQueryBody, err := json.Marshal(elasticQuery)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{\"error\": \"Fail to build elastic request.\"}")
		return
	}
	// Query elastic for text match.
	elasticRes, err := http.Post(rec.ElasticURL+"/search", "application/json", bytes.NewBuffer(elasticQueryBody))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{\"error\": \"Elastic request call failed with error "+err.Error()+"\"}")
		return
	}
	elasticResBytes, err := ioutil.ReadAll(elasticRes.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{\"error\": \"Cannot read elastic response.\"}")
		return
	}
	if elasticRes.StatusCode != 200 {
		w.WriteHeader(500)
		fmt.Fprintf(w,
			"{\"error\": \"Elastic request call rejeted with status "+
				elasticRes.Status+
				". Response is\n"+string(elasticResBytes)+"\n\"}")
	}

	var elasticResParsed elasticResponse
	err = json.Unmarshal(elasticResBytes, &elasticResParsed)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{\"error\": \"Cannot parse elastic response.\"}")
		return
	}

	// Collect target movies.
	requestOrder := make([]string, len(elasticResParsed.Hits.Hits))
	for idx, hit := range elasticResParsed.Hits.Hits {
		requestOrder[idx] = hit.Source.Name
	}
	instances := rec.buildInstancesFromMovieList(requestOrder)
	// Prepare TFX request.
	request := make(map[string][][2]float32)
	request["instances"] = instances
	requestBody, err := json.Marshal(request)
	if err != nil {
		fmt.Fprintf(w, "Error: tfx request failed due to "+err.Error())
		w.WriteHeader(500)
		return
	}
	// Make TFX calls.
	ranked, err := http.Post(rec.TfxURL+"/predict", "Content-Type: application/json; charset=utf-8", bytes.NewBuffer(requestBody))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: TFX request failed due to "+err.Error())
		return
	}
	body, _ := ioutil.ReadAll(ranked.Body)
	if ranked.StatusCode != 200 {
		w.WriteHeader(500)
		fmt.Fprintf(w, "ERROR: TFX call failed for "+string(body))
		return
	}
	// Load TFX response.
	rankResult := make(map[string][][]float32)
	json.Unmarshal(body, &rankResult)
	rankedMovies := make([]movie.Result, len(rankResult["predictions"]))
	for idx, score := range rankResult["predictions"] {
		rankedMovies[idx].Name = requestOrder[idx]
		rankedMovies[idx].Score = score[0]
	}
	// Sore ranked movies.
	sort.Slice(rankedMovies, func(i, j int) bool {
		return rankedMovies[i].Score < rankedMovies[j].Score
	})
	responseContent, err := json.Marshal(rankedMovies)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error: Serialize ranked result failed due to "+err.Error())
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseContent)
}

func (rec *recommender) handler() http.Handler {
	handler := http.NewServeMux()
	handler.HandleFunc("/movies", rec.movieRanking)
	handler.HandleFunc("/autocomplete", rec.autocomplete)
	return handler
}

func (rec *recommender) initMetadata(metadataFilePath string) error {
	// Load movie metadata from file.
	metadataFile, err := os.Open(metadataFilePath)
	defer metadataFile.Close()
	if err != nil {
		return err
	}
	r := csv.NewReader(bufio.NewReader(metadataFile))
	err = movie.BuildMoviesMetadata(r, &rec.Movies)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Parse flags to get backend urls.
	tfxURL := flag.String("tfx_url", "", "The url for tfx backend.")
	elasticURL := flag.String("elastic_url", "", "The url for elastic search service.")
	metadataFilePath := flag.String("metadata_path", "/data/movies_metadata.csv", "The path to the movie metadata csv file.")
	flag.Parse()

	// Create a recommender instance.
	rec := recommender{
		Movies:     make(map[string]movie.MetaData),
		TfxURL:     *tfxURL,
		ElasticURL: *elasticURL,
	}

	// Initialize the recommender.
	err := rec.initMetadata(*metadataFilePath)
	if err != nil {
		os.Stderr.WriteString("Load movie metadata failed.")
		return
	}
	// Create hanlders.
	http.ListenAndServe(":80", rec.handler())
}
