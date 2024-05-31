package db

import (
	"context"
	_ "embed"
	"github.com/jackc/pgx/v5/pgxpool"
	"shopTestTask/cfg"
)

type DB struct {
	*pgxpool.Pool
}

//go:embed schema.sql
var schema string

func New() DB {
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(context.Background(), schema)
	if err != nil {
		panic(err)
	}
	return DB{db}
}
