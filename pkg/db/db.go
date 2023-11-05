package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"soulmonk/dumper/pkg/db/ideas"
)

type Dao struct {
	conn         *pgx.Conn
	ctx          context.Context
	IdeasQuerier *ideas.Queries
}

func InitConnection(ctx context.Context, connectionString string) *pgx.Conn {
	var err error

	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		// todo no panic )
		panic(err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		// todo no panic )
		panic(err)
	}
	return conn
}

func GetDao(ctx context.Context, connectionString string) *Dao {
	dao := Dao{}
	dao.initConnection(ctx, connectionString)
	dao.initModels()

	return &dao
}

func (pg *Dao) Close() error {
	return pg.conn.Close(pg.ctx)
}

func (pg *Dao) GetDb() *pgx.Conn {
	return pg.conn
}

func (pg *Dao) initConnection(ctx context.Context, connectionString string) {
	pg.ctx = ctx
	pg.conn = InitConnection(ctx, connectionString)

	fmt.Println("Successfully connected!")
}

func (pg *Dao) initModels() {
	pg.IdeasQuerier = ideas.New(pg.conn)
}
