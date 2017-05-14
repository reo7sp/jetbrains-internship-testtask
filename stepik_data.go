package main

import (
	"time"
	"github.com/go-redis/redis"
)

type StepikLesson struct {
	Id uint64
	UpdateDate uint64
	TextStepIds []uint64
	AllStepIds []uint64
}

func FetchStepikLesson(pk uint64, redisClient *redis.Client) (StepikLesson, error) {
	var lesson StepikLesson
	lesson.Id = pk

	err := fetchBasicInfoAndSet(&lesson)
	if err != nil {
		return lesson, err
	}

	err = fetchStepIdsAndSet(&lesson, redisClient)
	if err != nil {
		return lesson, err
	}

	return lesson, nil
}

func fetchBasicInfoAndSet(lesson *StepikLesson) error {
	// fetch lesson from api
	apiResp, err := FetchStepikApiLesson(lesson.Id)
	if err != nil {
		return err
	}

	// additional parsing
	timeLayout := "2006-01-02T15:04:05Z"
	updateDate, err := time.Parse(timeLayout, apiResp.Lessons[0].UpdateDate)
	if err != nil {
		return err
	}
	updateDateUnix := uint64(updateDate.Unix())

	// construct object
	lesson.UpdateDate = updateDateUnix
	lesson.AllStepIds = apiResp.Lessons[0].Steps

	return nil
}

func fetchStepIdsAndSet(lesson *StepikLesson, redisClient *redis.Client) error {
	// get cache update time
	cacheUpdateDateUnix, err := FetchCacheLessonUpdateTime(redisClient, lesson.Id)
	if err != nil {
		cacheUpdateDateUnix = 0
	}

	// fetch step ids
	isCacheValid := lesson.UpdateDate <= cacheUpdateDateUnix
	if isCacheValid {
		err := fetchStepIdsFromCacheAndSet(lesson, redisClient)
		if err != nil {
			return err
		}
	} else {
		err := fetchStepIdsFromApiAndSet(lesson)
		if err != nil {
			return err
		}

		err = saveInCache(lesson, redisClient)
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchStepIdsFromCacheAndSet(lesson *StepikLesson, redisClient *redis.Client) error {
	stepIds, err := FetchCacheLessonTextStepsIds(redisClient, lesson.Id)
	if err != nil {
		return err
	}

	lesson.TextStepIds = stepIds

	return nil
}

func fetchStepIdsFromApiAndSet(lesson *StepikLesson) error {
	for _, stepId := range lesson.AllStepIds {
		isText, err := isStepText(stepId)
		if err != nil {
			return err
		}

		if isText {
			lesson.TextStepIds = append(lesson.TextStepIds, stepId)
		}
	}
	return nil
}

func isStepText(stepId uint64) (bool, error) {
	apiResp, err := FetchStepikApiStep(stepId)
	if err != nil {
		return false, err
	}

	isText := apiResp.Steps[0].Block.Name == "text"

	return isText, nil
}

func saveInCache(lesson *StepikLesson, redisClient *redis.Client) error {
	err := SaveCacheLessonTextStepsIds(redisClient, lesson.Id, lesson.TextStepIds)
	if err != nil {
		return err
	}

	err = SaveCacheLessonUpdateTime(redisClient, lesson.Id, lesson.UpdateDate)
	if err != nil {
		return err
	}

	return nil
}
