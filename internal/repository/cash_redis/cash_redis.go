package cashredis

import "github.com/go-redis/redis"

type CashRedis struct {
	Cash *redis.Client
}

type CashBannerActions interface {
	CashGetUserBanner() (struct{}, error)
}

func NewCashRedis() (*CashRedis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &CashRedis{
		Cash: rdb,
	}, nil
}
