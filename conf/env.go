package conf

import (
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
)

type Env struct {
	TMDBApiKey string `env:"TMDB_API_KEY"`
}

func LoadEnv() (Env, error) {
	_ = godotenv.Load(".env")
	var cfg Env
	if err := env.Parse(&cfg); err != nil {
		return Env{}, err
	}
	return cfg, nil
}
