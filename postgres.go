package yam

import (
	"database/sql"

	_ "github.com/lib/pq" // pg driver
)

type postgres struct {
	db *sql.DB
}

// assert postgres implements migrator
var _ migrator = &postgres{}

// NewPostgres creates a new postgres migrator
func NewPostgres(url string) (*postgres, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &postgres{db}, err
}

func (p *postgres) setup() error {
	if err := p.db.Ping(); err != nil {
		return err
	}

	var count int
	query := `SELECT COUNT(1) FROM information_schema.tables WHERE table_name = $1 AND table_schema = (SELECT current_schema()) LIMIT 1`
	if err := p.db.QueryRow(query, "data_migrations").Scan(&count); err != nil {
		return err
	}
	if count == 1 {
		return nil
	}

	query = `CREATE TABLE data_migrations (version bigint not null primary key)`
	if _, err := p.db.Exec(query); err != nil {
		return err
	}

	return nil
}

func (p *postgres) checkVersion(version int64) (bool, error) {
	var throwAway int64
	if err := p.db.QueryRow(`
		SELECT * from data_migrations
		WHERE version = $1
		LIMIT 1;
	`, version).Scan(&throwAway); err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (p *postgres) writeVersion(version int64) error {
	if _, err := p.db.Exec("INSERT INTO data_migrations VALUES ($1);", version); err != nil {
		return err
	}
	return nil
}

func (p *postgres) deleteVersion(version int64) error {
	if _, err := p.db.Exec("DELETE FROM data_migrations where version = $1;", version); err != nil {
		return err
	}
	return nil
}

func (p *postgres) teardown() error { return p.db.Close() }
