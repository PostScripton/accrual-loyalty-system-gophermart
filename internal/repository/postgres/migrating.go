package postgres

import (
	"context"
	"github.com/rs/zerolog/log"
	"os"
)

type migrationModel struct {
	id        int
	migration string
}

func (db *Postgres) MigrateAllMigrations(ctx context.Context) error {
	dir := "./migrations/"
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Debug().Err(err).Msg("Nothing to migrate")
		return err
	}

	runMigrations, err := db.getRunMigrations(ctx)
	if err != nil {
		return err
	}

	hasRun := 0
	for _, file := range files {
		migrationName := file.Name()[:len(file.Name())-4]

		_, ok := runMigrations[migrationName]
		if ok {
			continue
		}

		sql, err := os.ReadFile(dir + file.Name())
		if err != nil {
			log.Warn().Err(err).Msgf("Cannot open file [%s]", file.Name())
			return err
		}

		log.Debug().Msgf("[%s] Migrating...", migrationName)
		if err := db.migrate(ctx, migrationName, string(sql)); err != nil {
			log.Warn().Err(err).Msgf("[%s] Migration failed", migrationName)
			return err
		}
		hasRun++
		log.Debug().Msgf("[%s] Migrated!", migrationName)
	}

	if hasRun > 0 {
		log.Info().Msg("All migrations have run successfully!")
	} else {
		log.Info().Msg("Nothing to migrate")
	}
	return nil
}

func (db *Postgres) migrate(ctx context.Context, migrationName, migrationSQL string) error {
	if _, err := db.Exec(ctx, migrationSQL); err != nil {
		return err
	}

	saveMigrationSQL := `INSERT INTO migrations (migration) VALUES ($1);`
	if _, err := db.Exec(ctx, saveMigrationSQL, migrationName); err != nil {
		return err
	}

	return nil
}

func (db *Postgres) createMigrationsTable(ctx context.Context) error {
	sql := `
		CREATE TABLE IF NOT EXISTS migrations (
		    id			SERIAL PRIMARY KEY,
		    migration	VARCHAR(255) NOT NULL UNIQUE
		)
	`

	if _, err := db.Exec(ctx, sql); err != nil {
		return err
	}

	return nil
}

func (db *Postgres) getRunMigrations(ctx context.Context) (map[string]migrationModel, error) {
	if err := db.createMigrationsTable(ctx); err != nil {
		return map[string]migrationModel{}, err
	}

	runMigrations := make(map[string]migrationModel, 0)

	rows, err := db.Query(ctx, `SELECT id, migration FROM migrations;`)
	if err != nil {
		return map[string]migrationModel{}, err
	}

	for rows.Next() {
		runMigration := new(migrationModel)

		if err = rows.Scan(&runMigration.id, &runMigration.migration); err != nil {
			return map[string]migrationModel{}, err
		}

		runMigrations[runMigration.migration] = *runMigration
	}

	return runMigrations, nil
}
