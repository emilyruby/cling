package handlers 

import (
	"fmt"
	"net/http"
)

// HomeHandler handles requests to the root endpoint.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello world")
}