package repository

import (
	cashredis "banner/internal/repository/cash_redis"
	postgres "banner/internal/repository/postgres"
)

type Repository struct {
	Cash *cashredis.CashRedis
	DB   *postgres.Postgres
}

func HaveRepository() (*Repository, error) {
	cashRedis, err := cashredis.NewCashRedis()
	if err != nil {
		return nil, err
	}

	postgres, err := postgres.NewPostgres()
	if err != nil {
		return nil, err
	}

	return &Repository{
		Cash: cashRedis,
		DB:   postgres,
	}, nil
}
