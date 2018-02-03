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
	db, err := sql.Open("postgres", "postgres://localhost:5432/database?sslmode=enable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	migrations := []yam.Migration{
		{
			Version: 0,
			Up: func() error {
				_, err := db.Exec("INSERT INTO users (name) VALUES ($1)", "gopher")
				return err
			},
			Down: func() error {
				_, err := db.Exec("DELETE FROM users WHERE name = $1", "gopher")
				return err
			},
		},
	}

	yam.Sow(migrations, 0)
	yam.Reap(migrations, 0)
}

```

## Use cases
* Seeding data for different environments.
* Bootstrapping environments with data
* Running onetime scripts for data migrations
