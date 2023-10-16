package main

import (
	"fmt"
	"forum/internal/database"
	"net/http"
)

func main() {

	err := database.Initialize()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the forum!")
	})

	http.ListenAndServe(":8080", nil)
}
