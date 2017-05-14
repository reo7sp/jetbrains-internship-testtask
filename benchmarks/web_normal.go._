package main

import (
	"net/http"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"strconv"
	"time"
	"encoding/json"
)

func InitWeb(redisClient *redis.Client) (error) {
	partialHandleLessons := func(w http.ResponseWriter, r *http.Request) {
		handleLessons(redisClient, w, r)
	}
	http.HandleFunc("/lessons/", partialHandleLessons)

	err := http.ListenAndServe(":8000", nil)
	return err
}

func printError(w http.ResponseWriter, err error, tag string, status int) {
	fmt.Printf("error: %s: %s\n", tag, err)
	w.WriteHeader(status)
	fmt.Fprintf(w, "error: %s: %s", tag, err)
}

func handleLessons(redisClient *redis.Client, w http.ResponseWriter, r *http.Request) {
	path := strings.SplitN(r.URL.Path, "/", 3)
	pk, err := strconv.ParseUint(path[2], 10, 64)
	if err != nil {
		printError(w, err, "path parse", http.StatusNotFound)
		return
	}

	apiResp, err := FetchStepikApiLesson(pk)
	if err != nil {
		printError(w, err, "stepik api: fetch lesson", http.StatusNotFound)
		return
	}

	timeLayout := "2006-01-02T15:04:05Z"
	updateDate, err := time.Parse(timeLayout, apiResp.Lessons[0].UpdateDate)
	if err != nil {
		printError(w, err, "parse time", http.StatusInternalServerError)
		return
	}
	updateDateUnix := uint64(updateDate.Unix())

	cacheUpdateDateUnix, err := FetchCacheLessonUpdateTime(redisClient, pk)
	if err != nil {
		fmt.Printf("warning: cache: fetch lesson update time: %s", http.StatusInternalServerError)
		cacheUpdateDateUnix = 0
	}

	isCacheValid := updateDateUnix <= cacheUpdateDateUnix
	if isCacheValid {
		stepIds, err := FetchCacheLessonTextStepsIds(redisClient, pk)
		if err != nil {
			printError(w, err, "cache: fetch lesson text steps", http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(stepIds)
		if err != nil {
			printError(w, err, "json: marshal cached step ids", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	} else {
		var textStepIds []uint64
		for _, stepId := range apiResp.Lessons[0].Steps {
			apiResp, err := FetchStepikApiStep(stepId)
			if err != nil {
				printError(w, err, "stepik api: fetch step", http.StatusInternalServerError)
				return
			}
			if apiResp.Steps[0].Block.Name == "text" {
				textStepIds = append(textStepIds, stepId)
			}
		}

		err := SaveCacheLessonTextStepsIds(redisClient, pk, textStepIds)
		if err != nil {
			printError(w, err, "cache: save lesson text steps", http.StatusInternalServerError)
			return
		}
		err = SaveCacheLessonUpdateTime(redisClient, pk, updateDateUnix)
		if err != nil {
			printError(w, err, "cache: save lesson update time", http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(textStepIds)
		if err != nil {
			printError(w, err, "json: marshal step ids", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}
