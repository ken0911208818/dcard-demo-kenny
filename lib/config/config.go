package config

import (
	"log"
	"os"
	"strconv"
)

func GetStr(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Panic(`Environmental variable [` + key + `] don't exists`)
	}
	return value
}

func GetInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		log.Panic(`Environmental variable [` + key + `] don't exists`)
	}
	//字串轉int
	output, err := strconv.Atoi(value)

	if err != nil {
		log.Panic(`Environmental variable [` + key + `] is not an integer`)
	}
	return output
}

func GetBytes(key string) []byte {
	value := os.Getenv(key)
	if value == "" {
		log.Panic(`Environmental variable [` + key + `] don't exists`)
	}
	return []byte(value)
}
