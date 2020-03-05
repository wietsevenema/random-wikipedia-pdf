package main

import (
	"encoding/json"
	"fmt"
	"github.com/lithammer/shortuuid/v3"
	"log"
	"net/http"
	"wikipdf/lib/config"
	"wikipdf/lib/tasks"
)

type InsertResult struct {
	JobID string `json:"jobID"`
}

func (service *Service) InsertHandler(
	writer http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(writer)
	jobID := shortuuid.New()

	serviceUrl, err := service.run.GetServiceUrl(
		config.Region(),
		config.JobHandlingService(),
	)
	if err != nil {
		handleError(
			writer,
			"Error getting service URL",
			err,
		)
		return
	}

	task := &tasks.Task{
		Queue:          config.QueueName(),
		TaskID:         jobID,
		URL:            fmt.Sprintf("%s/pdfs/%s", serviceUrl, jobID),
		Body:           nil,
		ServiceAccount: config.ServiceAccount(),
	}

	err = service.storage.WriteInProgress(jobID, task)
	if err != nil {
		handleError(
			writer,
			"Error updating state",
			err,
		)
		return
	}
	err = service.tasks.AddTask(request.Context(), task)

	if err != nil {
		handleError(
			writer,
			"Error sending task",
			err,
		)
		return
	}

	err = encoder.Encode(InsertResult{JobID: jobID})
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}

}
