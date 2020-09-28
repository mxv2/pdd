package main

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := sqlx.MustOpen("sqlite3", "./test.db")

	ctx := context.Background()
	_, err := sqlx.LoadFileContext(ctx, db, "init.sql")
	if err != nil {
		panic(err)
	}

	log.Println("Init db: success")
}
