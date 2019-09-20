package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// EmployersHandler handles requests to the /employers endpoint and proxies requests to the Glassdoor API
func EmployersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	searchTerm := r.URL.Query().Get("search_term")
	glassdoorPartnerID := os.Getenv("GLASSDOOR_PARTNER_ID")
	glassdoorPartnerKey := os.Getenv("GLASSDOOR_PARTNER_KEY")

	glassdoorRequest, err := http.NewRequest("GET", "http://api.glassdoor.com/api/api.htm", nil)
	glassdoorRequestQuery := glassdoorRequest.URL.Query()
	glassdoorRequestQuery.Add("t.p", glassdoorPartnerID)
	glassdoorRequestQuery.Add("t.k", glassdoorPartnerKey)
	glassdoorRequestQuery.Add("format", "json")
	glassdoorRequestQuery.Add("action", "employers")
	glassdoorRequestQuery.Add("q", string(searchTerm))
	glassdoorRequest.URL.RawQuery = glassdoorRequestQuery.Encode()

	client := &http.Client{}
	response, err := client.Do(glassdoorRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unexpected server error")
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unexpected server error")
		return
	}

	w.WriteHeader(response.StatusCode)
	fmt.Fprintf(w, string(body))
}

func main() {
	port := os.Getenv("PORT")

	router := mux.NewRouter()
	router.
		HandleFunc("/employers", EmployersHandler).
		Methods("GET").
		Headers("Content-Type", "application/json")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
