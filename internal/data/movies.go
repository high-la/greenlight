package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/High-la/greenlight/internal/validator"
	"github.com/lib/pq"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // Use the - directive to hide it from json output
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"` // Add the omitempty directive
	// Use the Runtime type instead of int 32. Note that the omitempty directive will
	// still work on this: if the Runtime field has the underlying value 0, then it will
	// be considered empty and omitted -- and the MarshalJSON() method we jus made
	// won't be called at all.
	Runtime Runtime  `json:"runtime,omitempty"` // Add the omitempty directive
	Genres  []string `json:"genres,omitempty"`  // Add the omitempty directive
	Version int32    `json:"version"`
	// The version number starts at 1 and will be incremented each
	// time the movie information is updated
}

// Validate
func ValidateMovie(v *validator.Validator, movie *Movie) {

	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

// Define a MovieModel struct type which wraps a sql.DB connection pool.
type MovieModel struct {
	DB *sql.DB
}

// The Insert() method accepts a pointer to a movie struct, which should contain
// the data for the new record.
func (m MovieModel) Insert(movie *Movie) error {

	query := `
		INSERT INTO movies 
			(title, year, runtime, genres)
		VALUES 
			($1, $2, $3, $4)
		RETURNING 
			id, created_at, version`

	// Create an args slice containing the values for the placeholder paras from
	// the movie struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are bring used where* in the query
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic para and scanning the system
	// generated id, created_at and version values into the movie struct.
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)

}

func (m MovieModel) Get(id int64) (*Movie, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Define the SQL query for retrieving the movie data
	query := `
		SELECT 
			id, created_at, title, year, runtime, genres, version
		FROM movies
		WHERE id = $1`

	// Declare a Movie struct to hold data returned by the query
	var movie Movie

	// .
	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	// Handle any errors. If there was no matching movie found, scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Otherwise, return a pointer to the Movie struct
	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {

	// .
	query := `
		UPDATE movies
			SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		WHERE id = $5
		RETURNING version`

	// Create an args slice containing the values for the placeholder parameters.
	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	// .
	return m.DB.QueryRow(query, args...).Scan(&movie.Version)

}

// Add a placeholder method for deleting a specific record from the movies table.
func (m MovieModel) Delete(id int64) error {
	return nil
}
