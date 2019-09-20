package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")

	router := mux.NewRouter()
	router.
		HandleFunc("/employers", EmployersHandler).
		Methods("GET").
		Headers("Content-Type", "application/json")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
