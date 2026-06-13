package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gobasics/models"
)

// AlbumRepository defines the contract:
// what operations can be performed on albums 
// in the database
// This is interface-based design, 
// which allows for flexibility
// handler depends on this interface, 
// not on the concrete PostgresAlbumRepository 
// This promotes loose coupling and makes 
// it easier to test the handler,

type AlbumRepository interface {
	GetAll(ctx context.Context) ([]models.Album, error)
	GetByID(ctx context.Context, id int) (*models.Album, error)
	Create(ctx context.Context, req models.CreateAlbumRequest) (*models.Album, error)
	Update(ctx context.Context, id int , req models.CreateAlbumRequest) (*models.Album, error)
	Delete(ctx context.Context, id int) error
}

// PostgresAlbumRepository is a concrete 
// implementation of AlbumRepository
// that uses PostgreSQL as the data store. 
// It holds a reference to the sqlx.DB 
// connection pool which it uses 
// to execute queries against the database.

type PostgresAlbumRepository struct {
	db *sqlx.DB
}

func NewAlbumRepository(db *sqlx.DB) AlbumRepository {
	return &PostgresAlbumRepository{db: db}
}

func (r *PostgresAlbumRepository) GetAll(ctx context.Context) ([]models.Album, error) {
	var albums []models.Album
	// sqlx.Select scans multiple rows into a 
	// slice of structs, no manual row iterations
	err := r.db.SelectContext(ctx, &albums, 
		`SELECT id, title, price, artist_id, 
		created_at	 
		FROM albums
		ORDER BY created_at DESC
	`)

	return albums, err
}

// GetByID retrieves a single album by its ID.
// It uses sqlx.Get to scan a single row into a struct.
// If no rows are found, it returns a custom error.
func (r *PostgresAlbumRepository) GetByID(ctx context.Context, id int) (*models.Album, error) {
	var album models.Album
	err := r.db.GetContext(ctx, &album, 
		`SELECT id, title, price, artist_id, created_at 
		FROM albums 
		WHERE id = $1`, id)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("album with id %d not found", id)
		}
	return &album, err
}

// Create inserts a new album into the database 
// and returns the created album with its ID.
// Returning is a PostgreSQL feature that allows us to get
// the inserted row without a separate query.
func (r *PostgresAlbumRepository) Create(ctx context.Context,
	 req models.CreateAlbumRequest) (*models.Album, error) {
	var album models.Album
	err := r.db.QueryRowxContext(ctx, 
		`INSERT INTO albums (title, price, artist_id) 
		VALUES ($1, $2, $3) 
		RETURNING id, title, price, artist_id, created_at`, 
		req.Title, req.Price, req.ArtistID).StructScan(&album)
	return &album, err
}

// Update modifies an existing album's details in the database.
// It returns the updated album. 
// If the album does not exist, it returns a custom error.
func (r *PostgresAlbumRepository) Update(ctx context.Context, 
	id int, req models.CreateAlbumRequest)(*models.Album, error) {
	var album models.Album
	err := r.db.QueryRowxContext(ctx, 
		`UPDATE albums 
		SET title = $1, price = $2, artist_id = $3 
		WHERE id = $4 
		RETURNING id, title, price, artist_id, created_at`, 
		req.Title, req.Price, req.ArtistID, id).StructScan(&album)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("album with id %d not found", id)
	}
	return &album, err
}

// Delete removes an album from the database by its ID.
// If the album does not exist, it returns a custom error.
func (r *PostgresAlbumRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM albums WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("album with id %d not found", id)
	}
	return err

	//RowsAffected can be used to check if any 
	// row was actually deleted
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("album with id %d not found", id)
	}
	return nil
}
