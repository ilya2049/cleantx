package pgsql

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func NewConnection() (*pgx.Conn, func(), error) {
	conn, err := pgx.Connect(
		context.Background(),
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			"127.0.0.1",
			5432,
			"clinic-service",
			"password",
			"clinic-db",
		))
	if err != nil {
		return nil, nil, err
	}

	return conn, func() {
		if err := conn.Close(context.Background()); err != nil {
			log.Println("failed to close pg connection:", err)
		}
	}, nil
}
