package bot

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type Backend struct {
	Driver *sql.DB
	Config *BackendConf
}

type UserBack struct {
	hlogin  string
	command string
}

func NewBackend(back BackendConf) *Backend {
	return &Backend{
		Config: &back,
	}
}

func (b *Backend) Init() error {
	if err := b.Open(); err != nil {
		return err
	}
	defer b.Close()
	if _, err := b.Driver.Exec(`
	CREATE TABLE IF NOT EXISTS backend (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		hlogin VARCHAR(50), 
		command VARCHAR(50));`); err != nil {
		return err
	}
	if _, err := b.Driver.Exec(`
	CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		name VARCHAR(50),
		url VARCHAR(50), 
		data BLOB);`); err != nil {
		return err
	}
	return nil
}

func (b *Backend) Open() error {
	if b.Config.Dev {
		db, err := sql.Open("sqlite3", b.Config.Connection)
		if err != nil {
			return err
		}
		b.Driver = db
		return nil
	}
	return errors.New("database error")
}

func (b *Backend) Close() error {
	return b.Driver.Close()
}

func (b *Backend) GetLastCommand(hlogin string) (string, error) {
	if err := b.Open(); err != nil {
		return "", err
	}
	defer b.Close()
	query, err := b.Driver.Query("SELECT command FROM backend WHERE hlogin = $1 LIMIT 1", hlogin)
	if err != nil {
		return "", err
	}
	var command string
	for query.Next() {
		if err := query.Scan(&command); err != nil {
			return "", err
		}
	}
	if err := query.Close(); err != nil {
		return "", err
	}

	if _, err := b.Driver.Exec("DELETE FROM backend WHERE hlogin = $1", hlogin); err != nil {
		return "", err
	}
	return command, nil
}

func (b *Backend) DelLastCommands(hlogin string) error {
	if err := b.Open(); err != nil {
		return err
	}
	defer b.Close()
	_, err := b.Driver.Exec("DELETE FROM backend WHERE hlogin = $1", hlogin)
	if err != nil {
		return err
	}

	return nil
}

func (b *Backend) PutCommand(hlogin, command string) error {
	if err := b.Open(); err != nil {
		return err
	}
	defer b.Close()
	_, err := b.Driver.Exec(
		"INSERT INTO backend (hlogin, command) VALUES ($1, $2)",
		hlogin,
		command,
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *Backend) PutJson(data []byte, url string) error {
	if err := b.Open(); err != nil {
		return err
	}
	defer b.Close()
	_, err := b.Driver.Exec(
		"UPDATE questions SET data = $1 WHERE url = $2",
		data,
		url,
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *Backend) GetJson(url string) ([]byte, error) {
	if err := b.Open(); err != nil {
		return []byte{}, err
	}
	defer b.Close()
	query, err := b.Driver.Query("SELECT data FROM questions WHERE url = $1", url)
	if err != nil {
		return []byte{}, err
	}
	var data []byte
	for query.Next() {
		if err := query.Scan(&data); err != nil {
			return []byte{}, err
		}
	}
	if err := query.Close(); err != nil {
		return []byte{}, err
	}
	return data, err
}

func (b *Backend) GetPageUrls() ([]string, error) {
	if err := b.Open(); err != nil {
		return []string{}, err
	}
	defer b.Close()
	query, err := b.Driver.Query("SELECT url FROM questions")
	if err != nil {
		return []string{}, err
	}
	var data []string
	for query.Next() {
		var url string
		if err := query.Scan(&url); err != nil {
			return []string{}, err
		}
		data = append(data, url)
	}
	if err := query.Close(); err != nil {
		return []string{}, err
	}
	return data, err
}
