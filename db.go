package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PG interface {
	Open() error
	Close() error
	Insert(ctx context.Context, i int) error
}

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "test"
)

type pg struct {
	db *sql.DB
}

func NewPG() PG {
	return &pg{}
}

func (p *pg) Open() error {
	conn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return err
	}
	p.db = db
	if err = db.Ping(); err != nil {
		return err
	}
	fmt.Println("Opened DB connection")
	createQuery := `CREATE TABLE IF NOT EXISTS mytable (i integer);`
	_, err = db.Exec(createQuery)
	return err
}

func (p *pg) Close() (err error) {
	if p.db != nil {
		err = p.db.Close()
		fmt.Println("Closed DB connection")
		p.db = nil
	}
	return
}

func (p *pg) Insert(ctx context.Context, i int) error {
	insertStmt := fmt.Sprintf(`INSERT INTO "mytable"("i") values(%d)`, i)
	_, err := p.db.ExecContext(ctx, insertStmt)
	return err
}
