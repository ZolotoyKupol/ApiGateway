package migrations

import (
	"embed"
	"github.com/pkg/errors"
	"log/slog"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var migrations embed.FS

func Migrate(url string) error {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return errors.Wrap(err, "cannot connect to db")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return errors.Wrap(err, "cannot ping db")
	}

	entries,err := migrations.ReadDir(".")
	if err != nil {
		return errors.Wrap(err, "cannot read migrations")
	}
	slog.Info("found migrations", "count", len(entries))
	for _, entry := range entries {
		slog.Debug("migration", "name", entry.Name())
	}

	goose.SetBaseFS(migrations)
	if err = goose.SetDialect("postgres"); err != nil {
		return errors.Wrap(err, "cannot set dialect")
	}

	version, err := goose.GetDBVersion(db)
	if err != nil {
		return errors.Wrap(err, "cannot get version")
	}

	err = goose.Up(db, ".")
	if err != nil {
		if err = goose.DownTo(db, "migrations", version); err != nil {
			slog.Error("cannot rollback", slog.Any("error", err), slog.Any("version", version))
		}
		return errors.Wrap(err, "cannot migrate")
	}

	return nil
}
