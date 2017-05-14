package main

import (
	"net/http"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"strconv"
	"encoding/json"
)

func InitWeb(port uint16, redisClient *redis.Client) (error) {
	partialHandleLessons := func(w http.ResponseWriter, r *http.Request) {
		handleLessons(redisClient, w, r)
	}
	http.HandleFunc("/lessons/", partialHandleLessons)

	addr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(addr, nil)
	return err
}

func handleLessons(redisClient *redis.Client, w http.ResponseWriter, r *http.Request) {
	// path parse
	pk, err := extractPkFromPath(r.URL.Path)
	if err != nil {
		printError(w, err)
		return
	}

	// fetch lesson
	lesson, err := FetchStepikLesson(pk, redisClient)
	if err != nil {
		printError(w, err)
		return
	}

	// send response
	err = sendStepIdsAsJson(w, lesson.TextStepIds)
	if err != nil {
		printError(w, err)
		return
	}
}

func printError(w http.ResponseWriter, err error) {
	fmt.Printf("error: %s\n", err)
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "error: %s", err)
}

func extractPkFromPath(path string) (uint64, error) {
	pathParts := strings.SplitN(path, "/", 3)
	pk, err := strconv.ParseUint(pathParts[2], 10, 64)
	if err != nil {
		return pk, err
	}
	return pk, nil
}

func sendStepIdsAsJson(w http.ResponseWriter, stepIds []uint64) error {
	resp, err := json.Marshal(stepIds)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
	return nil
}
