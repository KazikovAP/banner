package redis

import "github.com/go-redis/redis"

type Redis struct {
	Cash *redis.Client
}

type CashBannerActions interface {
	CashGetUserBanner() (struct{}, error)
}

func NewCashRedis() (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Redis{
		Cash: rdb,
	}, nil
}
