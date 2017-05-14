package main

import (
	"log"
	"os"
	"strconv"
)

func getConfigOrDie() (string, string, int, uint16) {
	// redis addr
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// redis password
	redisPassword := os.Getenv("REDIS_ADDR")

	// redis db
	redisDbStr := os.Getenv("REDIS_DB")
	if redisDbStr == "" {
		redisDbStr = "0"
	}
	redisDbInt64, err := strconv.ParseInt(redisDbStr, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	redisDb := int(redisDbInt64)

	// web port
	webPortStr := os.Getenv("PORT")
	if webPortStr == "" {
		webPortStr = "8000"
	}
	webPortUint64, err := strconv.ParseUint(webPortStr, 10, 16)
	if err != nil {
		log.Fatal(err)
	}
	webPort := uint16(webPortUint64)

	return redisAddr, redisPassword, redisDb, webPort
}

func main() {
	redisAddr, redisPassword, redisDb, webPort := getConfigOrDie()

	redisClient, err := InitCache(redisAddr, redisPassword, redisDb)
	if err != nil {
		log.Fatal(err)
	}

	err = InitWeb(webPort, redisClient)
	if err != nil {
		log.Fatal(err)
	}
}
