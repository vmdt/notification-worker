package redis

type RedisOptions struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Password      string `mapstructure:"password"`
	Database      int    `mapstructure:"database"`
	PoolSize      int    `mapstructure:"pool_size"`
	EnableTracing bool   `mapstructure:"enable_tracing"`
	Uri           string `mapstructure:"uri"`
}
