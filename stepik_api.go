package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io/ioutil"
)

type StepikApiLesson struct {
	Id         uint64 `json:"id"`
	Steps      []uint64 `json:"steps"`
	UpdateDate string `json:"update_date"`
}

type StepikApiLessonsResponse struct {
	Lessons []StepikApiLesson `json:"lessons"`
}

func FetchStepikApiLesson(pk uint64) (lesson StepikApiLessonsResponse, err error) {
	resp, err := http.Get(fmt.Sprintf("https://stepik.org:443/api/lessons/%d", pk))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &lesson)
	return
}

type StepikApiStepBlock struct {
	Name string `json:"name"`
}

type StepikApiStep struct {
	Id    uint64 `json:"id"`
	Block StepikApiStepBlock `json:"block"`
}

type StepikApiStepsResponse struct {
	Steps []StepikApiStep `json:"steps"`
}

func FetchStepikApiStep(pk uint64) (step StepikApiStepsResponse, err error) {
	resp, err := http.Get(fmt.Sprintf("https://stepik.org:443/api/steps/%d", pk))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &step)
	return
}
