package repo

import (
	"banner/internal/lib/logger"
	"banner/internal/models"
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewUser(db *pgxpool.Pool, log *slog.Logger) *UserRepo {
	return &UserRepo{db, log}
}

func (u *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	_, err := u.db.Exec(ctx,
		`
		INSERT INTO users (username, password, role)
		VALUES ($1, $2, $3)
		RETURNING id
		`, user.Username, user.Password, user.Role)

	if err != nil {
		u.log.Error("Failed to create user", logger.Err(err))
		return err
	}

	return nil
}

func (u *UserRepo) FindUserUsername(ctx context.Context, username string) (models.User, error) {
	query, err := u.db.Query(ctx, `SELECT * FROM users WHERE username = $1`, username)
	if err != nil {
		u.log.Error("Error querying users", logger.Err(err))
		return models.User{}, err
	}
	defer query.Close()

	order := models.User{}
	if !query.Next() {
		u.log.Error("User not found")
		return models.User{}, fmt.Errorf("User not found")
	} else {
		err := query.Scan(&order.ID, &order.Username, &order.Password, &order.Role)
		if err != nil {
			u.log.Error("Error scanning users", logger.Err(err))
			return models.User{}, err
		}
	}

	return order, nil
}

func (u *UserRepo) FindUserId(ctx context.Context, id int) (models.User, error) {
	query, err := u.db.Query(ctx, `SELECT * FROM users WHERE id = $1`, id)
	if err != nil {
		u.log.Error("Error querying users", logger.Err(err))
		return models.User{}, err
	}
	defer query.Close()

	order := models.User{}
	if !query.Next() {
		u.log.Error("User not found")
		return models.User{}, fmt.Errorf("User not found")
	} else {
		err := query.Scan(&order.ID, &order.Username, &order.Password, &order.Role)
		if err != nil {
			u.log.Error("Error scanning users", logger.Err(err))
			return models.User{}, err
		}
	}

	return order, nil
}
