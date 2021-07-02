package config

import (
	"fmt"
	"os"
)

func GetConfigInf() (string, string) {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PWD")
	database := os.Getenv("MYSQL_AUTH_DB")
	return "mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user,
		password, host, port, database)
}
