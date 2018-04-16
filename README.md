[![CircleCI](https://circleci.com/gh/cbelsole/yam.svg?style=svg)](https://circleci.com/gh/cbelsole/yam)

# YAM (yet another migrator)

There are [plenty of migrators](https://awesome-go.com/#database) out there most of which deal with db schema migrations in some form or another. YAM is another migrator that aims to be light weight in capabilities but heavy in usefulness.

## Example:
```go
package main

import (
	"database/sql"

	"github.com/cbelsole/yam"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	migrations := []yam.Migration{
		{
			Version: 0,
			Up: func() error {
				_, err := db.Exec(`CREATE TABLE users(
					id SERIAL PRIMARY KEY,
					name text not null
				);`)
				return err
			},
			Down: func() error {
				_, err := db.Exec("DROP TABLE USERS")
				return err
			},
		},
		{
			Version: 1,
			Up: func() error {
				_, err := db.Exec("INSERT INTO users (name) VALUES ($1);", "gopher")
				return err
			},
			Down: func() error {
				_, err := db.Exec("DELETE FROM users WHERE name = $1;", "gopher")
				return err
			},
		},
	}

	// Running migrations without migrator skips version checks
	if err := yam.Migrate(nil, migrations, 0); err != nil {
		panic(err)
	}
	if err := yam.Rollback(nil, migrations, 0); err != nil {
		panic(err)
	}

	pg, err := yam.NewPostgres("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = yam.Migrate(pg, migrations, 0); err != nil {
		panic(err)
	}

	pg, err = yam.NewPostgres("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = yam.Rollback(pg, migrations, 0); err != nil {
		panic(err)
	}

	// NewPostgresFromDB does not clean up the connection. So you can reuse the
	// migrator.
	pg = yam.NewPostgresFromDB(db)
	if err = yam.Migrate(pg, migrations, 0); err != nil {
		panic(err)
	}

	if err = yam.Rollback(pg, migrations, 0); err != nil {
		panic(err)
	}
}
```

## Use cases
* Seeding data for different environments.
* Bootstrapping environments with data
* Running onetime scripts for data migrations

## Running tests
```
# testing without integrations
go test ./... -short

# testing with integrations
docker-compose up
go test ./...
```

## Integrations
* postgres

Pull requests welcome
