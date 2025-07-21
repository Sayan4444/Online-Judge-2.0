package utils

import (
	"os"

	"github.com/joho/godotenv"
)

const configPath = "../.env"

func LoadEnv() {
	err := godotenv.Load(configPath)
	if err != nil {
		panic(err)
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
