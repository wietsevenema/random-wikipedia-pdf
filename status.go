package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type State string

const (
	StateInProgress State = "IN_PROGRESS"
	StateCompleted        = "COMPLETED"
)

type StateResult struct {
	JobID string `json:"jobID"`
	State State  `json:"state"`
}

func (service *Service) StatusHandler(
	writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	jobID := vars["jobID"]
	state := StateResult{
		JobID: jobID,
	}
	encoder := json.NewEncoder(writer)

	result, err := service.storage.Completed(
		request.Context(),
		jobID,
	)
	if err != nil {
		handleError(writer, "Error reading state", err)
		return
	}
	if result {
		state.State = StateCompleted
		err := encoder.Encode(state)
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
		return
	}

	inProgress, err := service.storage.InProgress(request.Context(), jobID)
	if err != nil {
		handleError(writer, "Error reading state", err)
		return
	}
	if inProgress {
		state.State = StateInProgress

		err = encoder.Encode(state)
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
		return
	}

	writer.WriteHeader(http.StatusNotFound)
	return

}
