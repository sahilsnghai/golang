package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var EnvsInit = initConfig()

type Config struct {
	PublicHost            string
	Port                  string
	DBUser                string
	DBPassword            string
	DBAddress             string
	DBName                string
	SecretKey             string
	JWTExpirationInSecond int64
}

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	return Config{
		PublicHost:            getEnv("MYSQL_HOST", "localhost"),
		Port:                  getEnv("MYSQL_PORT", "3306"),
		DBUser:                getEnv("MYSQL_USER", "root"),
		DBPassword:            getEnv("MYSQL_PASSWORD", "mypassword"),
		DBAddress:             fmt.Sprintf("%s:%s", getEnv("MYSQL_HOST", "localhost"), getEnv("MYSQL_PORT", "3306")),
		DBName:                getEnv("MYSQL_DB", "Project6"),
		SecretKey:             getEnv("secretKey", "adbakshbdaisvdlsdjavbsd"),
		JWTExpirationInSecond: getEnvInt("JWT_EXP", 3600*24*7),
	}

}

func getEnv(key, fallback string) string {

	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvInt(key string, fallback int64) int64 {

	if value, ok := os.LookupEnv(key); ok {
		value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return value
	}
	return fallback
}
