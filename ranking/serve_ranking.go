package main


import (
	"fmt"
	"net/http"
)


func main() {
	fmt.Println("Hello, World!")
}


func handler(w http.ResponseWriter, r *http.Request) {

    fmt.Fprintf(w, "Movie List %s,%s,%s,%s", 
	            "Nope","Thor: Love and Thunder","Minions: The Rise of Gru",
				"Where the Crawdads Sing")
}