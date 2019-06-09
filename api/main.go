package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/emilyruby/cling/api/handlers"
)

func main() {
	r := mux.NewRouter()
  r.HandleFunc("/", handlers.HomeHandler)
	log.Fatal(http.ListenAndServe(":8000", r))
}