package yam

import (
	"sort"
	"strconv"
)

// Migration interface describes the way migrations can be migrated or rolled
// back
type (
	Migration struct {
		Version int
		Up      func() error
		Down    func() error
	}
	migrationSlice []Migration
)

func (m Migration) String() string { return strconv.Itoa(m.Version) }

func (m migrationSlice) Len() int           { return len(m) }
func (m migrationSlice) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m migrationSlice) Less(i, j int) bool { return m[i].Version < m[j].Version }
func (m migrationSlice) first(n int) migrationSlice {
	if n == 0 {
		return m
	}
	return m[0 : len(m)-(len(m)-n)]
}

// Sow runs migrations without validating with additional offset to run the
// first (offset) migrations
func Sow(migrations migrationSlice, offset int) error {
	sort.Sort(migrations)
	for _, migration := range migrations.first(offset) {
		if migration.Up != nil {
			if err := migration.Up(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Reap runs migration without validating in reverse order
func Reap(migrations migrationSlice, offset int) error {
	sort.Sort(sort.Reverse(migrations))
	for _, migration := range migrations.first(offset) {
		if migration.Down != nil {
			if err := migration.Down(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Migrate runs migrations after validating that they have not been run
func Migrate(dburl string, migrations migrationSlice, offset int) {

}

// Rollback runs migrations after validating that they have not been run
func Rollback(dburl string, migrations migrationSlice, offset int) {

}
