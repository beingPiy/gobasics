package repository

import (
	"context" // for managing request contexts and timeouts
	"database/sql" // for working with SQL databases
	"errors" // for creating and handling errors
	"fmt" // for formatting strings and errors
	
	"gobasics/models" // for the User model and related request/response structs
	"github.com/jmoiron/sqlx" // for working with SQL databases with additional features
)

type UserRepository interface {
	Create(ctx context.Context, email, passwordHash string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type PostgresUserRepository struct {
	db *sqlx.DB // database connection pool
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, email, passwordHash string) (*models.User, error) {
	query := `INSERT INTO users (email, password_hash) 
			  VALUES ($1, $2) 
			  RETURNING id, email, role, created_at`
	var user models.User
	err := r.db.QueryRowContext(ctx, query, email, passwordHash).Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt)
	
	if err != nil && err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
		return nil, fmt.Errorf("email already registered")
	}

	return &user, err
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, role, created_at 
			  FROM users
			  WHERE email = $1`
	var user models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user not found")
	}

	return &user, err
}



