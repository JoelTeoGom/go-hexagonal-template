package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // registers the "mysql" driver
)

// Open creates the connection pool and checks it with a ping, so connection
// problems show up at startup instead of on the first query. The caller owns
// the returned pool and must close it.
func Open(ctx context.Context, dsn string) (*sql.DB, error) {
	// sql.Open does not connect: it only validates the DSN and prepares the pool.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("mysql: open: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Ping does open a real connection, so the first error surfaces here.
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("mysql: ping: %w", err)
	}
	return db, nil
}
