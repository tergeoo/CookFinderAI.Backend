package db

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql" //mysql driver
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

// PostgresSQLX
// как пользоваться https://jmoiron.github.io/sqlx/
func PostgresSQLX(ds Datasource) (*sqlx.DB, error) {
	sources := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", ds.User, ds.Pass, ds.Host, ds.Port, ds.Name)
	db, err := sqlx.Connect("pgx", sources)
	if err != nil {
		return nil, fmt.Errorf("utils: failed to connect to db (%s) - %w", ds.Name, err)
	}

	_, err = db.Exec("SET TIME ZONE 'UTC'")
	if err != nil {
		return nil, fmt.Errorf("utils: failed to set time zone - %w", err)
	}

	_, err = db.Exec(fmt.Sprintf("SET search_path TO %s", ds.Schema))
	if err != nil {
		return nil, fmt.Errorf("utils: failed to set search path - %w", err)
	}

	return db, nil
}

func ClickhouseSQLX(ds Datasource) (*sqlx.DB, error) {
	source := fmt.Sprintf("clickhouse://%s:%s@%s:%d/%s", ds.User, ds.Pass, ds.Host, ds.Port, ds.Name)

	db, err := sqlx.Connect("clickhouse", source)
	if err != nil {
		return nil, fmt.Errorf("utils: failed to connect to ClickHouse (%s) - %w", ds.Name, err)
	}

	return db, nil
}

func MySQLSQLX(ds Datasource) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=UTC",
		ds.User, // Пользователь базы данных
		ds.Pass, // Пароль
		ds.Host, // Хост
		ds.Port, // Порт
		ds.Name, // Имя базы данных
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("utils: failed to connect to db (%s) - %w", ds.Name, err)
	}

	return db, nil
}

func Migrate(db *sql.DB, migrations embed.FS, appName, migrationsDir string) error {
	slog.Info("utils: start applying migrations", "embed", migrations)

	goose.SetBaseFS(migrations)

	appName = strings.ReplaceAll(appName, "-", "_")
	goose.SetTableName("migrations_" + appName)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("utils: failed to set goose dialect - %w", err)
	}

	if migrationsDir == "" {
		migrationsDir = "."
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		if errors.Is(err, goose.ErrNoNextVersion) {
			slog.Info("utils: no new migrations to apply")
			return nil
		}
		return fmt.Errorf("utils: failed to apply migrations - %w", err)
	}

	return nil
}

func MigrateClickHouse(db *sql.DB, migrations embed.FS, appName string) error {
	table := "migrations_" + strings.ReplaceAll(appName, "-", "_")

	_, err := db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			version String,
			applied_at DateTime DEFAULT now()
		) ENGINE = MergeTree()
		ORDER BY version;
	`, table))
	if err != nil {
		return fmt.Errorf("failed to create migration history table: %w", err)
	}

	applied := make(map[string]bool)
	rows, err := db.Query(fmt.Sprintf(`SELECT version FROM %s`, table))
	if err != nil {
		return fmt.Errorf("failed to read migration history: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[v] = true
	}

	entries, err := migrations.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read embedded migrations dir: %w", err)
	}

	var toApply []string
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".sql") || entry.IsDir() {
			continue
		}
		if !applied[name] {
			toApply = append(toApply, name)
		}
	}

	sort.Strings(toApply)

	for _, file := range toApply {
		sqlBytes, err := migrations.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}

		if _, err := db.Exec(string(sqlBytes)); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", file, err)
		}

		if _, err := db.Exec(fmt.Sprintf(`INSERT INTO %s (version) VALUES (?)`, table), file); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}

		fmt.Println("✔ Applied:", file)
	}

	return nil
}

func PostgresPGX(ds Datasource) (*pgxpool.Pool, error) {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		ds.User, ds.Pass, ds.Host, ds.Port, ds.Name,
	)
	db, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(context.Background(), "SET TIME ZONE 'UTC'")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(context.Background(), fmt.Sprintf("SET search_path TO %s", ds.Schema))
	if err != nil {
		return nil, err
	}

	return db, nil
}
