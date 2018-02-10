package yam

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"

	_ "github.com/lib/pq" // pg driver
)

func TestPostgresMigrate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping postgres tests")
	}
	data := make(map[string]int, 0)
	tests := []struct {
		name       string
		migrations migrationSlice
		migrator   func(*testing.T) migrator
		offset     int64
		err        error
		expected   int
	}{
		{
			name:     "migratePostgres",
			migrator: testPostgresMigrator,
			migrations: []Migration{
				{
					Version: 1,
					Up: func() error {
						if data["migratePostgres"] != 0 {
							return errors.New(`data["migratePostgres"] should be 0`)
						}
						data["migratePostgres"]++
						return nil
					},
				},
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Migrate(tt.migrator(t), tt.migrations, tt.offset); err != nil && err != tt.err {
				t.Errorf("Migrate() expected `%v` but received `%v`", tt.err, err)
			} else if data[tt.name] != tt.expected {
				t.Errorf("Migrate() expected `%d` but received `%d`", tt.expected, data[tt.name])
			}

			var versions []int64
			for _, migration := range tt.migrations {
				if migration.Up != nil {
					versions = append(versions, migration.Version)
				}
			}

			gotVersions := getVersions()

			if !reflect.DeepEqual(gotVersions, versions) {
				t.Errorf("Migration versions want `%v` got `%v`", versions, gotVersions)
			}

			cleanupMigrations()
		})
	}
}

func TestPostgresRollback(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping postgres tests")
	}

	data := make(map[string]int, 0)
	tests := []struct {
		name       string
		migrations migrationSlice
		migrator   func(*testing.T) migrator
		offset     int64
		err        error
		expected   int
	}{
		{
			name:     "migrateRollback",
			migrator: testPostgresMigrator,
			migrations: []Migration{
				{
					Version: 1,
					Up:      func() error { return nil },
					Down: func() error {
						if data["migrateRollback"] != 0 {
							return errors.New(`data["migrateRollback"] should be 0`)
						}
						data["migrateRollback"]++
						return nil
					},
				},
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.name)
			testMigrator := testPostgresMigrator(t)
			if err := Migrate(testMigrator, tt.migrations, 0); err != nil {
				t.Errorf("test migrator error: %s", err)
			}

			if err := Rollback(tt.migrator(t), tt.migrations, tt.offset); err != nil && err != tt.err {
				t.Errorf("Rollback() expected `%v` but received `%v`", tt.err, err)
			} else if data[tt.name] != tt.expected {
				t.Errorf("Rollback() expected `%d` but received `%d`", tt.expected, data[tt.name])
			}

			if gotVersions := getVersions(); len(gotVersions) != 0 {
				t.Errorf("Migration versions len should be `0` got `%v`", len(gotVersions))
			}
		})
	}
}

func getVersions() []int64 {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var versions []int64
	rows, err := db.Query("SELECT * from data_migrations")
	if err != nil {
		if err == sql.ErrNoRows {
			return versions
		}
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			panic(err)
		}
		versions = append(versions, version)
	}
	return versions
}

func cleanupMigrations() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if _, err := db.Exec("TRUNCATE data_migrations"); err != nil {
		panic(err)
	}

}

func testPostgresMigrator(t *testing.T) migrator {
	pg, err := NewPostgres("postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		t.Error("error initializing pg migrator for `migratePostgres`")
	}
	return pg

}
