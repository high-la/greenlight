package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Add a createMovieHandler for the "POST /v1/moview" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint. For now retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	// When httprouter is parsing a request, any interpolated URL params will be
	// stored in the request context. We can use the ParamsFromContext() fun to
	// retrieve a slice containing these para names and values.
	params := httprouter.ParamsFromContext(r.Context())

	// We can then use the ByName method to get the value of the params
	// from the slice. the value returned by ByName is always a string
	// So we try to convert it to a base 10 integer (with a bit size of 64).
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Otherwise, interpolate the movie ID in a placeholder response.
	fmt.Fprintf(w, "show the details of movie %d\n", id)
}
