package main

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"wikipdf/lib/tasks"
	"net/http"
)

type State string

const (
	StateInProgress State = "IN_PROGRESS"
	StateCompleted        = "COMPLETED"
)

type StateResult struct {
	JobID string `json:jobID`
	State State  `json:state`
}

func (service *Service) StateHandler(
	writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	jobID := vars["jobID"]
	state := StateResult{
		JobID: jobID,
	}
	encoder := json.NewEncoder(writer)

	if _, e := service.bucket.Object(jobID + ".pdf").Attrs(request.Context());
		e != storage.ErrObjectNotExist {
		state.State = StateCompleted
		encoder.Encode(state)
		writer.WriteHeader(http.StatusOK)
		return
	}

	if _, e := service.bucket.Object(jobID + ".pending").Attrs(request.Context());
		e != storage.ErrObjectNotExist {
		state.State = StateInProgress
		encoder.Encode(state)
		writer.WriteHeader(http.StatusOK)
		return
	}

	writer.WriteHeader(http.StatusNotFound)
	return

}

//FIXME extract this to separate package (.pdf, .pending, are spreading everywhere)
func (service *Service) WritePending(jobID string, task *tasks.Task) error {
	sw := service.bucket.
		Object(jobID + ".pending").
		NewWriter(context.Background())
	defer sw.Close()

	en := json.NewEncoder(sw)
	return en.Encode(task)
}
