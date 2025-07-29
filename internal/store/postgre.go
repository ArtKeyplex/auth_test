package store

import (
	"auth_test/configs"
	"auth_test/internal/service"
	"context"
	"github.com/jackc/pgx/v5"
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

func (ps *PostgresUserStore) Add(user *service.User) (bool, error) {

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
