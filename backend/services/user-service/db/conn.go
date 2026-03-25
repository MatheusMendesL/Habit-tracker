package db

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func Conn() (*sql.DB, *Queries, error) {
	_ = godotenv.Load()

	typeDB := os.Getenv("TYPEDB")
	dbString := os.Getenv("DBSTRING")

	db, err := sql.Open(typeDB, dbString)
	if err != nil {
		return nil, nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, nil, err
	}

	queries := New(db)

	return db, queries, nil
}
