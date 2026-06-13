package models

import "time"

// Artist mirrors the artists table in the database
// The `db` tags specify the column names for sqlx 
// to map struct fields to database columns

type Artist struct {
	ID        int       `db:"id"			json:"id"`
	Name      string    `db:"name"			json:"name"`
	Country   string    `db:"country"		json:"country"`
	CreatedAt time.Time `db:"created_at"	json:"created_at"`
}

// Album mirrors the albums table in the database
type Album struct {
	ID        int       `db:"id"			json:"id"`
	Title     string    `db:"title"		json:"title" validate:"required",min=3,max=100"`
	Price	 float64   `db:"price"		json:"price" validate:"required,gt=0"`
	ArtistID  int       `db:"artist_id"	json:"artist_id" validate:"required"`
	CreatedAt time.Time `db:"created_at"	json:"created_at"`
}

// CreateAlbumRequest is used for validating incoming JSON payloads
// when creating a new album. It includes validation tags 
// to ensure the data meets certain criteria before processing.
type CreateAlbumRequest struct {
	Title    string  `json:"title" validate:"required,min=3,max=100"`
	Price    float64 `json:"price" validate:"required,gt=0"`
	ArtistID int     `json:"artist_id" validate:"required"`
}