package config

type Config struct {
	SessionLimit int `env:"SESSION_LIMIT" envDefault:"10"`
}
