package bot

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Backend struct {
	DriverMaster *sql.DB
	ConfigMaster *BackendConf
	DriverSlave  *sql.DB
	ConfigSlave  *BackendConf
}

func NewBackend(backMaster, backSlave BackendConf) *Backend {
	backend := Backend{
		ConfigMaster: &backMaster,
	}
	if backMaster.Multi {
		backend.ConfigSlave = &backSlave
	} else {
		backend.ConfigSlave = &backMaster
	}

	dbHostR := os.Getenv("HOST_READ")
	dbHostW := os.Getenv("HOST_WRITE")
	if dbHostR != "" {
		backend.ConfigMaster.Host = dbHostR
	}
	if dbHostR != "" {
		backend.ConfigSlave.Host = dbHostW
	}
	return &backend
}

func (b *Backend) Init() error {
	if err := b.OpenMaster(); err != nil {
		return err
	}
	if err := CreateTables(b.DriverMaster, b.ConfigMaster.Dev); err != nil {
		return err
	}
	b.CloseMaster()
	// if b.ConfigMaster.Multi {
	// 	if err := b.OpenSlave(); err != nil {
	// 		return err
	// 	}
	// 	if err := CreateTables(b.DriverSlave, b.ConfigMaster.Dev); err != nil {
	// 		return err
	// 	}
	// 	b.CloseSlave()
	// }
	return nil
}

func CreateTables(driver *sql.DB, dev bool) error {
	if dev {
		if _, err := driver.Exec(`
	CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		name VARCHAR(50) UNIQUE,
		url VARCHAR(50) UNIQUE, 
		data BLOB);`); err != nil {
			return err
		}
		return nil
	} else {
		if _, err := driver.Exec(`
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

func (b *Backend) OpenMaster() error {
	if b.ConfigMaster.Dev {
		db, err := sql.Open("sqlite3", "./sqlite.db")
		if err != nil {
			return err
		}
		b.DriverMaster = db
		return nil
	} else {
		conn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			b.ConfigMaster.Host, b.ConfigMaster.Port, b.ConfigMaster.User,
			b.ConfigMaster.Password, b.ConfigMaster.DatabaseName, b.ConfigMaster.SSLMode,
		)

		db, err := sql.Open("postgres", conn)
		if err != nil {
			return err
		}
		b.DriverMaster = db
		return nil
	}
}

func (b *Backend) OpenSlave() error {
	if b.ConfigSlave.Dev {
		db, err := sql.Open("sqlite3", "./sqlite.db")
		if err != nil {
			return err
		}
		b.DriverSlave = db
		return nil
	} else {
		conn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			b.ConfigSlave.Host, b.ConfigSlave.Port, b.ConfigSlave.User,
			b.ConfigSlave.Password, b.ConfigSlave.DatabaseName, b.ConfigSlave.SSLMode,
		)

		db, err := sql.Open("postgres", conn)
		if err != nil {
			return err
		}
		b.DriverSlave = db
		return nil
	}
}

func (b *Backend) CloseMaster() error {
	return b.DriverMaster.Close()
}

func (b *Backend) CloseSlave() error {
	return b.DriverSlave.Close()
}

func (b *Backend) UpdatePage(data []byte, url string) error {
	if err := b.OpenMaster(); err != nil {
		return err
	}
	defer b.CloseMaster()
	_, err := b.DriverMaster.Exec(
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
	if err := b.OpenMaster(); err != nil {
		return err
	}
	defer b.CloseMaster()
	_, err := b.DriverMaster.Exec(
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
	if err := b.OpenSlave(); err != nil {
		return []byte{}, err
	}
	defer b.CloseSlave()
	query, err := b.DriverSlave.Query("SELECT data FROM questions WHERE url = $1", url)
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
	if err := b.OpenSlave(); err != nil {
		return []byte{}, err
	}
	defer b.CloseSlave()
	query, err := b.DriverSlave.Query("SELECT data FROM questions WHERE name = $1", name)
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
	if err := b.OpenSlave(); err != nil {
		return []string{}, []string{}, err
	}
	defer b.CloseSlave()
	query, err := b.DriverSlave.Query("SELECT name, url FROM questions")
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
