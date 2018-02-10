package yam

import (
	"errors"
	"testing"
)

func TestMigrate(t *testing.T) {
	migrationError := errors.New("migration error")
	data := make(map[string]int, 0)
	tests := []struct {
		name       string
		migrations migrationSlice
		offset     int64
		err        error
		expected   int
	}{
		{
			name: "migrate",
			migrations: []Migration{
				{
					Version: 1,
					Up: func() error {
						if data["migrate"] != 0 {
							return errors.New(`data["migrate"] should be 0`)
						}
						data["migrate"]++
						return nil
					},
				},
				{
					Version: 2,
					Up: func() error {
						if data["migrate"] != 1 {
							return errors.New(`data["migrate"] should be 1`)
						}
						data["migrate"]++
						return nil
					},
				},
				{
					Version: 3,
				},
			},
			expected: 2,
		},
		{
			name: "migrateWithOffset",
			migrations: []Migration{
				{
					Version: 1,
					Up: func() error {
						if data["migrateWithOffset"] != 0 {
							return errors.New(`data["migrateWithOffset"] should be 0`)
						}
						data["migrateWithOffset"]++
						return nil
					},
				},
				{
					Version: 2,
					Up: func() error {
						if data["migrateWithOffset"] != 1 {
							return errors.New(`data["migrateWithOffset"] should be 1`)
						}
						data["migrateWithOffset"]++
						return nil
					},
				},
				{
					Version: 3,
				},
			},
			offset:   1,
			expected: 1,
		},
		{
			name: "error",
			migrations: []Migration{
				{
					Version: 1,
					Up: func() error {
						return migrationError
					},
				},
			},
			err: migrationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Migrate(nil, tt.migrations, tt.offset); err != nil && err != tt.err {
				t.Errorf("Migrate() expected `%v` but received `%v`", tt.err, err)
			} else if data[tt.name] != tt.expected {
				t.Errorf("Migrate() expected `%d` but received `%d`", tt.expected, data[tt.name])
			}
		})
	}
}

func TestRollback(t *testing.T) {
	migrationError := errors.New("migration error")
	data := make(map[string]int, 0)
	tests := []struct {
		name       string
		migrations migrationSlice
		offset     int64
		err        error
		expected   int
	}{
		{
			name: "rollback",
			migrations: []Migration{
				{
					Version: 1,
					Down: func() error {
						if data["rollback"] != 1 {
							return errors.New(`data["rollback"] should be 1`)
						}
						data["rollback"]++
						return nil
					},
				},
				{
					Version: 2,
					Down: func() error {
						if data["rollback"] != 0 {
							return errors.New(`data["rollback"] should be 0`)
						}
						data["rollback"]++
						return nil
					},
				},
				{
					Version: 3,
				},
			},
			expected: 2,
		},
		{
			name: "rollbackWithOffset",
			migrations: []Migration{
				{
					Version: 1,
					Down: func() error {
						if data["rollbackWithOffset"] != 1 {
							return errors.New(`data["rollbackWithOffset"] should be 1`)
						}
						data["rollbackWithOffset"]++
						return nil
					},
				},
				{
					Version: 2,
					Down: func() error {
						if data["rollbackWithOffset"] != 0 {
							return errors.New(`data["rollbackWithOffset"] should be 0`)
						}
						data["rollbackWithOffset"]++
						return nil
					},
				},
				{
					Version: 3,
				},
			},
			offset:   2,
			expected: 1,
		},
		{
			name: "error",
			migrations: []Migration{
				{
					Version: 1,
					Down: func() error {
						return migrationError
					},
				},
			},
			err: migrationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Rollback(nil, tt.migrations, tt.offset); err != nil && err != tt.err {
				t.Errorf("Rollback() expected `%v` but received `%v`", tt.err, err)
			} else if data[tt.name] != tt.expected {
				t.Errorf("Rollback() expected `%d` but received `%d`", tt.expected, data[tt.name])
			}
		})
	}
}
