package initializers

import (
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"os"
)

type Env struct {
	Port              string `env:"PORT" envDefault:"8080"`
	LogFile           string `env:"LOG_FILE" envDefault:"gin.log"`
	IntroFilePath     string `env:"INTRO_FILE_PATH" envDefault:"app/assets/intro.mkv"`
	MovieSourceFolder string `env:"MOVIE_SOURCE_FOLDER" envDefault:"./"`
	MovieTargetFolder string `env:"MOVIE_TARGET_FOLDER" envDefault:"./"`
	TvSourceFolder    string `env:"TV_SOURCE_FOLDER" envDefault:"./"`
	TvTargetFolder    string `env:"TV_TARGET_FOLDER" envDefault:"./"`
	TMDBApiKey        string `env:"TMDB_API_KEY" envDefault:""`
	DBSync            bool   `env:"DB_SYNC" envDefault:"false"`
	DBHost            string `env:"DB_HOST" envDefault:"localhost"`
	DBPort            string `env:"DB_PORT" envDefault:"5432"`
	DBUser            string `env:"DB_USER" envDefault:"postgres"`
	DBPassword        string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName            string `env:"DB_NAME" envDefault:"postgres"`
	RedisHost         string `env:"REDIS_HOST" envDefault:"localhost:6379"`
	RedisPassword     string `env:"REDIS_PASSWORD" envDefault:""`
	S3AccessKeyId     string `env:"S3_ACCESS_KEY_ID" envDefault:""`
	S3SecretAccessKey string `env:"S3_SECRET_ACCESS_KEY" envDefault:""`
	S3BucketName      string `env:"S3_BUCKET_NAME" envDefault:""`
	S3Endpoint        string `env:"S3_ENDPOINT" envDefault:"https://s3.fr-par.scw.cloud"`
	ScanCron          string `env:"SCAN_CRON" envDefault:"*/15 * * * *"`
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
