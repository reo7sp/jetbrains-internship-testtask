package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisPassword := os.Getenv("REDIS_ADDR")
	redisDbStr := os.Getenv("REDIS_DB")
	if redisDbStr == "" {
		redisDbStr = "0"
	}
	redisDb, err := strconv.ParseInt(redisDbStr, 10, 32)
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := InitCache(redisAddr, redisPassword, int(redisDb))
	if err != nil {
		log.Fatal(err)
	}

	err = InitWeb(redisClient)
	if err != nil {
		log.Fatal(err)
	}
}
