package bot

import (
	"database/sql"
	"errors"
)

type Backend struct {
	Database *sql.DB
	Config   *ConfigDB
}

func NewBackend(back Backend) *Backend {
	return &Backend{
		Config: back.Config,
	}
}

func (b *Backend) Open() error {
	if b.Config.DebugON {
		db, err := sql.Open("sqlite3", ":memory:")
		if err != nil {
			return err
		}
		b.Database = db
		return nil
	}
	return errors.New("Database error")
}

func (b *Backend) Close() error {
	return b.Database.Close()
}
