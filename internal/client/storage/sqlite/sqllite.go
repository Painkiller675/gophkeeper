package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/Painkiller675/gophkeeper/internal/proto"
	"github.com/Painkiller675/gophkeeper/internal/server/storage"
)

// localSecret is used to return data from local database
type localSecret struct {
	Content []byte
	Version string
}

// LocalStorage represents a local client sqlite storage struct
type LocalStorage struct {
	LocalDB *sql.DB
}

// NewLocalStorage - a constructor of local sqlite database struct
func NewLocalStorage() (*LocalStorage, error) {
	lDB, err := initLocalDG()
	if err != nil {
		return nil, err
	}
	return &LocalStorage{
		LocalDB: lDB,
	}, nil
}

// initLocalDG initializes the SQLite database
func initLocalDG() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./local.db")
	if err != nil {
		return nil, err
	}

	// Create a table for client's local db

	createTableSQL := `CREATE TABLE IF NOT EXISTS secrets (
		id SERIAL PRIMARY KEY,
		name VARCHAN(255) NOT NULL,
    	content BYTEA,
    	version UUID NOT NULL UNIQUE
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// GetLocalSecret - returns a specific secret in offline mode
func (ls *LocalStorage) GetLocalSecret(ctx context.Context, name string) (*localSecret, error) {
	row := ls.LocalDB.QueryRowContext(
		ctx,
		`SELECT content, version FROM secrets WHERE name = ($1)`,
		name,
	)
	localSec := &localSecret{}
	err := row.Scan(&localSec.Content, &localSec.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return &localSecret{}, storage.ErrSecretNotFound
	}
	return localSec, err
}

// SyncLocalSecrets - synchronize remote database with a local one
func (ls *LocalStorage) SyncLocalSecrets(ctx context.Context, secret *proto.SecretInfo) error {
	// Update local SQLite
	_, err := ls.LocalDB.Exec("INSERT OR REPLACE INTO secrets (name, content, version) VALUES (?, ?, ?)",
		secret.GetName(), secret.GetContent(), secret.GetVersion())
	if err != nil {
		log.Println("Error saving to local database:", err)
	}
	return nil
}
