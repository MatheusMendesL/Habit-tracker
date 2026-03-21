package db

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Conn() (*sql.DB, *Queries, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/habit_tracker")
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
