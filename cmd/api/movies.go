package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/High-la/greenlight/internal/data"
	"github.com/High-la/greenlight/internal/validator"
)

// Add a createMovieHandler for the "POST /v1/moview" endpoint. For now we simply
// return a plain-text placeholder response.
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	// Declare an anonymous struct to hold the infn that we expect to be in the
	// HTTP request body (note that the field names and types in the struct are a subset
	// of the Moview struct that we created earlier). This struct will be our *target
	// decode destination*.
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	// Use the new readJSON() helper to decode the request body into the input struct.
	// If this returns an error we send the client the error along with a 400
	// Bad Request status code.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Initialize a new Validator instance.
	v := validator.New()

	// Use the Check() method to execute our validation checks. This will add the
	// provided key and error message to the errors map if the check doesn't evaluate
	// to true. For example, in the first line here we "check" that the title is not
	// equal to the empty string". In the second, we "check the length of the title
	// is less than or equal to 500 bytes" and so on.
	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "ust not be more than 500 bytes long")

	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genres")
	v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
	// Note that we're usinf the unique helper in the line below to check that all
	// values in the input.Geners slice are unique.
	v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")

	// Use the Valid() method to see if any of the checks failed. If they did, then use
	// the failedValidationResponse() helper to send a response to the client, passing
	// in the v.Errors map.
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Dump the contents of the  input struct in a HTTP response.
	fmt.Fprintf(w, "%+v\n", input)
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint. For now retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		// Use the new notFoundResponse() helper.
		app.notFoundResponse(w, r)
		return
	}

	// Create a new instance of the Movie struct, containing the ID we extracted from
	// the URL and some dummy data. Also  notice that we deliberately haven't set a
	// value for the Year field.
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Geners:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	// Encode the struct to JSON and send it as the HTTP response.

	// Create an envelope{"movie": movie} instance and pass it towriteJSON(), instead
	// of passing the plain movie struct.
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		// Use the new serverErrorResponse()
		app.serverErrorResponse(w, r, err)
	}

}
