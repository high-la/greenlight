package data

import "time"

type Movie struct {
	ID        int64
	CreatedAt time.Time // Timestamp for when the movie is added to our database
	Title     string
	Year      int32    // Movie release year
	Runtime   int32    // movie length in minutes
	Geners    []string // slice of geners
	Version   int32
	// The version number starts at 1 and will be incremented each
	// time the movie information is updated
}
