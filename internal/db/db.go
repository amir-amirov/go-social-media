package db

import (
	"database/sql"
	"log"
	"time"
)

func New(addr string, maxOpenConns, maxIdleConns int) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	for i := 0; i < 10; i++ {
		if err = db.Ping(); err == nil {
			break
		}
		time.Sleep(time.Second * 2)
	}

	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to database..")

	return db, nil
}
