package cache

type RedisConfig struct {
	RedisHost     string
	RedisPort     int
	RedisUsername string
	RedisPassword string
	RedisDB       int
	RedisUseSSL   bool
}

// NewRedisConfig 返回 Redis 配置。
func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		RedisHost:     "localhost",
		RedisPort:     6379,
		RedisUsername: "",
		RedisPassword: "",
		RedisDB:       0,
		RedisUseSSL:   false,
	}
}
