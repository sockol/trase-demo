package utils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

type IDB interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions, fn txFn) error
	Open() (*sql.DB, error)
}

type DB struct{}

func NewDB() DB {
	return DB{}
}

func (db *DB) Open() (*sql.DB, error) {
	return open("public")
}

type txFn func(tx *sql.Tx) error

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions, fn txFn) error {
	return beginTx(db, ctx, opts, fn)
}

type TestDB struct {
	Fixture
}

type Fixture struct {
	UserId1 uuid.UUID
	UserId2 uuid.UUID
	PostId1 uuid.UUID
	PostId2 uuid.UUID
}

func TestNewDB(t *testing.T) TestDB {
	// hardcode test creds because I want to run tests outside of docker through vscode
	t.Setenv("POSTGRES_USER", "postgres")
	t.Setenv("POSTGRES_PASSWORD", "posgres349")
	t.Setenv("POSTGRES_DB_NAME", "postgres")
	t.Setenv("POSTGRES_DB_HOST", "localhost")

	userId1, _ := uuid.Parse("4a2b9c10-9daf-11ed-93ce-0242ac120001")
	userId2, _ := uuid.Parse("4a2b9c10-9daf-11ed-93ce-0242ac120002")
	postId1, _ := uuid.Parse("4a2b9c10-9daf-11ed-93ce-0242ac220001")
	postId2, _ := uuid.Parse("4a2b9c10-9daf-11ed-93ce-0242ac220002")
	f := Fixture{
		UserId1: userId1,
		UserId2: userId2,
		PostId1: postId1,
		PostId2: postId2,
	}
	return TestDB{
		Fixture: f,
	}
}

// Open clears the test db and opens a connection.
func (db *TestDB) Open() (*sql.DB, error) {
	conn, err := open("test")
	if err != nil {
		return nil, err
	}
	sql := `
	DO $$ DECLARE
		r RECORD;
	BEGIN
		FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = current_schema()) LOOP
			EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' cascade';
		END LOOP;
	END $$;
	`
	_, err = conn.Exec(sql)
	if err != nil {
		return nil, err
	}

	// Fixtures
	_, err = conn.Query(`INSERT INTO users (id, name, email) VALUES ($1, 'user-1', 'email-1');`, db.Fixture.UserId1)
	if err != nil {
		return nil, err
	}
	_, err = conn.Query(`INSERT INTO users (id, name, email) VALUES ($1, 'user-2', 'email-2');`, db.Fixture.UserId2)
	if err != nil {
		return nil, err
	}
	_, err = conn.Query(`INSERT INTO posts (id, title, content, user_id) VALUES ($1, 'title-1', 'content-1', $2);`, db.Fixture.PostId1, db.Fixture.UserId1)
	if err != nil {
		return nil, err
	}
	_, err = conn.Query(`INSERT INTO posts (id,title, content, user_id) VALUES ($1, 'title-2', 'content-2', $2);`, db.Fixture.PostId2, db.Fixture.UserId2)
	if err != nil {
		return nil, err
	}

	s := ""
	err = conn.QueryRow(`select current_schema;`).Scan(&s)
	if err != nil {
		return nil, err
	}

	return conn, err
}

func (db *TestDB) BeginTx(ctx context.Context, opts *sql.TxOptions, fn txFn) error {
	return beginTx(db, ctx, opts, fn)
}

func beginTx(db IDB, ctx context.Context, opts *sql.TxOptions, fn txFn) error {
	conn, err := db.Open()
	if err != nil {
		return err
	}
	tx, err := conn.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	err = fn(tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func open(schema string) (*sql.DB, error) {
	connStr := ""
	// This will be set on prod by heroku.
	if url, ok := os.LookupEnv("DATABASE_URL"); ok {
		connStr = url
	} else {
		connStr = fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable search_path=%s",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB_NAME"),
			os.Getenv("POSTGRES_DB_HOST"),
			schema,
		)
	}
	return sql.Open("postgres", connStr)
}
