package yam

import (
	"errors"
	"testing"
)

func TestSow(t *testing.T) {
	migrationError := errors.New("migration error")
	data := make(map[string]int, 0)
	tests := []struct {
		name       string
		migrations MigrationSlice
		offset     int
		err        error
		expected   int
	}{
		{
			name: "sow",
			migrations: []Migration{
				{
					Version: 2,
					Up: func() error {
						if data["sow"] != 1 {
							return errors.New(`data["sow"] should be 1`)
						}
						data["sow"]++
						return nil
					},
				},
				{
					Version: 1,
					Up: func() error {
						if data["sow"] != 0 {
							return errors.New(`data["sow"] should be 0`)
						}
						data["sow"]++
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
			name: "sow with offset",
			migrations: []Migration{
				{
					Version: 2,
					Up: func() error {
						if data["sow with offset"] != 1 {
							return errors.New(`data["sow with offset"] should be 1`)
						}
						data["sow with offset"]++
						return nil
					},
				},
				{
					Version: 1,
					Up: func() error {
						if data["sow with offset"] != 0 {
							return errors.New(`data["sow with offset"] should be 0`)
						}
						data["sow with offset"]++
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
			if err := Sow(tt.migrations, tt.offset); err != nil && err != tt.err {
				t.Errorf("Sow() expected `%v` but received `%v`", tt.err, err)
			} else if data[tt.name] != tt.expected {
				t.Errorf("Sow() expected `%d` but received `%d`", tt.expected, data[tt.name])
			}
		})
	}
}

func TestReap(t *testing.T) {
	migrationError := errors.New("migration error")
	data := make(map[string]int, 0)
	tests := []struct {
		name       string
		migrations MigrationSlice
		offset     int
		err        error
		expected   int
	}{
		{
			name: "reap",
			migrations: []Migration{
				{
					Version: 2,
					Down: func() error {
						if data["reap"] != 0 {
							return errors.New(`data["reap"] should be 0`)
						}
						data["reap"]++
						return nil
					},
				},
				{
					Version: 1,
					Down: func() error {
						if data["reap"] != 1 {
							return errors.New(`data["reap"] should be 1`)
						}
						data["reap"]++
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
			name: "reap with offset",
			migrations: []Migration{
				{
					Version: 2,
					Down: func() error {
						if data["reap with offset"] != 0 {
							return errors.New(`data["reap with offset"] should be 0`)
						}
						data["reap with offset"]++
						return nil
					},
				},
				{
					Version: 1,
					Down: func() error {
						if data["reap with offset"] != 1 {
							return errors.New(`data["reap with offset"] should be 1`)
						}
						data["reap with offset"]++
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
			if err := Reap(tt.migrations, tt.offset); err != nil && err != tt.err {
				t.Errorf("Reap() expected `%v` but received `%v`", tt.err, err)
			} else if data[tt.name] != tt.expected {
				t.Errorf("Reap() expected `%d` but received `%d`", tt.expected, data[tt.name])
			}
		})
	}
}
