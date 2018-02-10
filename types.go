package yam

type (
	// Migration describes how to migrate and roll back changes with a version to
	// keep track of changes
	Migration struct {
		Version int64
		Up      func() error
		Down    func() error
	}
	migrationSlice []Migration
	migrator       interface {
		checkVersion(int64) (bool, error)
		deleteVersion(int64) error
		setup() error
		teardown() error
		writeVersion(int64) error
	}
)

func (m migrationSlice) first(n int64) migrationSlice {
	if n == 0 {
		return m
	}
	return m[0 : int64(len(m))-(int64(len(m))-n)]
}
