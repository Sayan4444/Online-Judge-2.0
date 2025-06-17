/***********************************************************************
     Copyright (c) 2025 GNU/Linux Users' Group (NIT Durgapur)
************************************************************************/

package config

import (
	// "fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// path for env file
const configPath = "./.env"

// Load environment variables from env file
func LoadEnv() {
	err := godotenv.Load(configPath)
	if err != nil {
		panic(err)
	}
}

// retrieve env value from key
func GetEnv(key string) string {
	val := os.Getenv(key)

	return val
}

// Connect to database
func ConnectDB() (*gorm.DB, error) {
	dsn := GetEnv("DSN_STRING")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}
