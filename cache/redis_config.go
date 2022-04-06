package cache

import (
	"os"
	"strconv"
)

type RedisConfig struct {
	Host               string
	Port               string
	Password           string
	DB                 int
	CacheTTLSeconds    int
	ShouldDisableCache bool
}

func (redisConfig *RedisConfig) Init() {
	redisConfig.ShouldDisableCache, _ = strconv.ParseBool(os.Getenv("REDIS_SHOULD_DISABLE_CACHE"))
	redisConfig.Host = os.Getenv("REDIS_HOST")
	redisConfig.Port = os.Getenv("REDIS_PORT")
	redisConfig.Password = os.Getenv("REDIS_PASSWORD")
	redisConfig.DB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))
	redisConfig.CacheTTLSeconds, _ = strconv.Atoi(os.Getenv("REDIS_CACHE_TTL_SECONDS"))
}
