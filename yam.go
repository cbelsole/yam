package yam

// Migrate runs migrations after validating that they have not been run
func Migrate(m Migrator, migrations migrationSlice, offset int64) error {
	if m != nil {
		if err := m.setup(); err != nil {
			return err
		}
	}

	for _, migration := range migrations.first(offset) {
		var skip bool
		if m != nil {
			migrated, err := m.checkVersion(migration.Version)
			if err != nil {
				return err
			}
			skip = !migrated
		}

		if migration.Up != nil && !skip {
			if err := migration.Up(); err != nil {
				return err
			}

			if m != nil {
				if err := m.writeVersion(migration.Version); err != nil {
					return err
				}
			}
		}
	}

	if m != nil {
		if err := m.teardown(); err != nil {
			return err
		}
	}

	return nil
}

// Rollback runs migrations after validating that they have not been run
func Rollback(m Migrator, migrations migrationSlice, offset int64) error {
	if m != nil {
		if err := m.setup(); err != nil {
			return err
		}
	}

	rollbacks := make(migrationSlice, len(migrations))
	for i := 0; i < len(migrations); i++ {
		rollbacks[i] = migrations[len(migrations)-i-1]
	}

	for _, rollback := range rollbacks.first(offset) {
		var skip bool
		if m != nil {
			migrated, err := m.checkVersion(rollback.Version)
			if err != nil {
				return err
			}
			skip = migrated
		}

		if rollback.Down != nil && !skip {
			if err := rollback.Down(); err != nil {
				return err
			}

			if m != nil {
				if err := m.deleteVersion(rollback.Version); err != nil {
					return err
				}
			}
		}
	}

	if m != nil {
		if err := m.teardown(); err != nil {
			return err
		}
	}
	return nil
}
