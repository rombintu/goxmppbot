package bot

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Backend struct {
	Driver *sql.DB
	Config *BackendConf
}

type UserBack struct {
	HLogin  string
	Command string
}

func NewBackend(back BackendConf) *Backend {
	backend := Backend{
		Config: &back,
	}
	dbName := os.Getenv("DATABASE")
	if dbName != "" {
		backend.Config.DatabaseName = dbName
		backend.Config.Host = os.Getenv("HOST")
		backend.Config.Port = os.Getenv("PORT")
		backend.Config.User = os.Getenv("USER")
		backend.Config.Password = os.Getenv("PASSWORD")
		backend.Config.SSLMode = os.Getenv("SSLMODE")
	}
	return &backend
}

func (b *Backend) Init() error {
	if err := b.Open(); err != nil {
		return err
	}
	defer b.Close()
	if b.Config.Dev {
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
		name VARCHAR(50) UNIQUE,
		url VARCHAR(50) UNIQUE, 
		data BLOB);`); err != nil {
			return err
		}
		return nil
	} else {
		if _, err := b.Driver.Exec(`
	CREATE TABLE IF NOT EXISTS backend (
		id SERIAL, 
		hlogin VARCHAR(50), 
		command VARCHAR(50));`); err != nil {
			return err
		}
		if _, err := b.Driver.Exec(`
	CREATE TABLE IF NOT EXISTS questions (
		id SERIAL, 
		name VARCHAR(50) UNIQUE,
		url VARCHAR(50) UNIQUE, 
		data BYTEA);`); err != nil {
			return err
		}
		return nil
	}
}

func (b *Backend) Open() error {
	if b.Config.Dev {
		db, err := sql.Open("sqlite3", "./sqlite.db")
		if err != nil {
			return err
		}
		b.Driver = db
		return nil
	} else {
		conn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			b.Config.Host, b.Config.Port, b.Config.User, b.Config.Password, b.Config.DatabaseName, b.Config.SSLMode,
		)

		db, err := sql.Open("postgres", conn)
		if err != nil {
			return err
		}
		b.Driver = db
		return nil
	}
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

func (b *Backend) UpdatePage(data []byte, url string) error {
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

func (b *Backend) PutNewPage(name, url string) error {
	if err := b.Open(); err != nil {
		return err
	}
	defer b.Close()
	_, err := b.Driver.Exec(
		"INSERT INTO questions (name, url) VALUES ($1, $2)",
		name,
		url,
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *Backend) GetJsonByUrl(url string) ([]byte, error) {
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

func (b *Backend) GetJsonByName(name string) ([]byte, error) {
	if err := b.Open(); err != nil {
		return []byte{}, err
	}
	defer b.Close()
	query, err := b.Driver.Query("SELECT data FROM questions WHERE name = $1", name)
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

func (b *Backend) GetPageUrlsAndNames() ([]string, []string, error) {
	if err := b.Open(); err != nil {
		return []string{}, []string{}, err
	}
	defer b.Close()
	query, err := b.Driver.Query("SELECT name, url FROM questions")
	if err != nil {
		return []string{}, []string{}, err
	}
	var names []string
	var urls []string
	for query.Next() {
		var name, url string
		if err := query.Scan(&name, &url); err != nil {
			return []string{}, []string{}, err
		}
		names = append(names, name)
		urls = append(urls, url)
	}
	if err := query.Close(); err != nil {
		return []string{}, []string{}, err
	}
	return urls, names, err
}
