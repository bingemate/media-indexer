package initializers

import (
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"os"
)

type Env struct {
	Port              string `env:"PORT" envDefault:"8080"`
	LogFile           string `env:"LOG_FILE" envDefault:"gin.log"`
	MovieSourceFolder string `env:"MOVIE_SOURCE_FOLDER" envDefault:"./"`
	MovieTargetFolder string `env:"MOVIE_TARGET_FOLDER" envDefault:"./"`
	TvSourceFolder    string `env:"TV_SOURCE_FOLDER" envDefault:"./"`
	TvTargetFolder    string `env:"TV_TARGET_FOLDER" envDefault:"./"`
	TMDBApiKey        string `env:"TMDB_API_KEY" envDefault:""`
}

func LoadEnv() (Env, error) {
	var envCfg = &Env{}

	err := godotenv.Load(".env")
	if err != nil && !os.IsNotExist(err) {
		return Env{}, err
	}

	if err := env.Parse(envCfg); err != nil {
		return Env{}, err
	}
	return *envCfg, nil
}
