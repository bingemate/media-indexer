package initializers

import (
	"fmt"
	"github.com/bingemate/media-go-pkg/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func ConnectToDB(env Env) (*gorm.DB, error) {
	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Paris",
		env.DBHost, env.DBUser, env.DBPassword, env.DBName, env.DBPort)
	log.Println("Connecting to database...")
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("Connected to database")

	if env.DBSync {
		log.Println("Syncing database...")
		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
		err = repository.Migrate(db)
		if err != nil {
			return nil, err
		}
		log.Println("Database synced")
	}
	return db, nil
}
