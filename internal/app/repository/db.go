package repository

import (
	"context"
	"database/sql"
	"github.com/blagorodov/go-shortener/internal/app/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

var Connection *sql.DB

func InitDB() {
	db, err := sql.Open("pgx", config.Options.DBDataSource)
	if err != nil {
		panic(err)
	}
	Connection = db
}

func CloseDB() {
	Connection.Close()
}

func PingDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return Connection.PingContext(ctx)
}
