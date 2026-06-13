package db

import (
	"log"	
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // register PostgreSQL driver - blank
)

// Connect opens a connection pool to the database 
// using the provided connection string and verifies the connection.
// Why a connection pool, not a single connection? 
// A pool keeps N connections open and manages them, 
// by handing them to goroutines on demand, 
// and closing them when not needed. 
// This is more efficient than opening and closing a 
// single connection for each request, which would cost
// TCP + TLS handshakes cost on every request(~10ms).
// sqlx.DB manages a pool of connections, 
// allowing for concurrent access and better performance 
// in a web application context.
func Connect(connectionString string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Tuning the connection pool settings for 
	// better performance
	// MaxOpenConns: The maximum number of open connections 
	// to the database.
	// MaxIdleConns: The maximum number of connections 
	// in the idle connection pool.
	// ConnMaxLifetime: The maximum amount of time a 
	// connection may be reused.
	db.SetMaxOpenConns(25) // Adjust based on your workload and database limits
	db.SetMaxIdleConns(25) // Keep some idle connections for quick reuse
	db.SetConnMaxLifetime(5 * 60) // 5 
	
	log.Println("Successfully connected to PostgreSQL")
	return db
}