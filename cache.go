package main

import (
	"github.com/go-redis/redis"
	"fmt"
	"strconv"
)

func InitCache(redisAddr string, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
		DB:       db,
	})

	pong, err := client.Ping().Result()
	if err != nil || pong != "PONG" {
		return nil, err
	}
	return client, err
}

func FetchCacheLessonUpdateTime(client *redis.Client, pk uint64) (uint64, error) {
	key := fmt.Sprintf("lesson:%d:update-date", pk)
	value, err := client.Get(key).Uint64()
	return value, err
}

func SaveCacheLessonUpdateTime(client *redis.Client, pk uint64, updateTime uint64) error {
	key := fmt.Sprintf("lesson:%d:update-date", pk)
	value := strconv.FormatUint(updateTime, 10)
	err := client.Set(key, value, 0).Err()
	return err
}

func FetchCacheLessonTextStepsIds(client *redis.Client, pk uint64) ([]uint64, error) {
	key := fmt.Sprintf("lesson:%d:text-steps", pk)
	members, err := client.SMembers(key).Result()
	if err != nil {
		return nil, err
	}

	stepIds := make([]uint64, len(members))
	for i, member := range members {
		id, err := strconv.ParseUint(member, 10, 64)
		if err != nil {
			return nil, err
		}
		stepIds[i] = id
	}
	return stepIds, nil
}

func SaveCacheLessonTextStepsIds(client *redis.Client, pk uint64, stepIds []uint64) error {
	key := fmt.Sprintf("lesson:%d:text-steps", pk)
	client.Del(key)
	for _, id := range stepIds {
		value := strconv.FormatUint(id, 10)
		err := client.SAdd(key, value).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
