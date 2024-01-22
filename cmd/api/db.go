package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// sqlを開く
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// sql接続
func (app *application) connectToDB() (*sql.DB, error) {
	connection, err := openDB(app.DSN)
	if err != nil {
		return nil, err
	}

	log.Println("Connection to Postgres")
	return connection, nil

}
