package store

import (
	"auth_test/configs"
	"auth_test/internal/service"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"regexp"
)

type PostgresUserStore struct {
	db *pgx.Conn
}

func NewPostgresUserStore(db *pgx.Conn) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

func (ps *PostgresUserStore) Get(username string) (*service.User, error) {
	var user service.User
	err := ps.db.QueryRow(context.Background(), "SELECT login, password FROM users WHERE login=$1", username).
		Scan(&user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ps *PostgresUserStore) AddUser(ctx context.Context, db *pgx.Conn, login string, password string) (bool, error) {
	log.Info().Msg("adding user")
	_, err := db.Exec(ctx,
		`INSERT INTO users (login, password) VALUES ($1, $2)`,
		login, password)
	if err != nil {
		err = FindError(err.Error())
		log.Error().Err(err).Msg("error adding user")
		return false, err
	}
	log.Info().Msg("successfully added user")
	return true, nil
}

func InitDb(ctx context.Context, cfg *configs.Config) (*pgx.Conn, error) {
	db, err := pgx.Connect(ctx, cfg.DSN)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(ctx,
		`
		CREATE TABLE IF NOT EXISTS users (
			id serial PRIMARY KEY,
			login TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);
	`)
	if err != nil {
		db.Close(ctx)
		return nil, err
	}

	return db, nil
}

func FindError(errMsg string) error {
	re := regexp.MustCompile(`SQLSTATE (\d{5})`)
	errCode := re.FindStringSubmatch(errMsg)[1]

	log.Info().Str("Code", errCode).Msg("finding error")
	switch errCode {
	case "23505":
		return service.ErrUserExists
	case "23502":
		return service.NullNotAllowedError
	default:
		return service.SomethingWrongError
	}
}
