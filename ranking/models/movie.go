package movie

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// MetaData holds metadata for moveis.
type MetaData struct {
	Name        string
	VoteCount   float32
	VoteAverage float32
}

// Result holds the movie name and its ranking score.
type Result struct {
	Name  string
	Score float32
}

// BuildMoviesMetadata loads movie metadata from a csv file and returnt them in a
// map.
func BuildMoviesMetadata(r *csv.Reader, movies *map[string]MetaData) error {
	// Skip the header line.
	_, err := r.Read()
	if (err != nil) {
		fmt.Println("Failed to loading csv header for ", err.Error())
		return err
	}
	line := 1
	for {
		line++ 
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Reading csv failed at line " + strconv.Itoa(line)+ " for error ", err.Error())
			continue
		}
		voteCount, err := strconv.ParseFloat(record[23], 32)
		if err != nil {
			fmt.Println("Loading record vote count failed for ", err.Error(), "\nRecord is ", record)
			continue
		}
		voteAverage, err := strconv.ParseFloat(record[22], 32)
		if err != nil {
			fmt.Println("Loading record vote average failed for ", err.Error(), "\nRecord is ", record)
			continue
		}
		(*movies)[record[8]] = MetaData{
			Name:        record[8],
			VoteCount:   float32(voteCount),
			VoteAverage: float32(voteAverage),
		}
	}
	return nil
}