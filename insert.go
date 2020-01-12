package main

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/ksuid"
	"wikipdf/lib/config"
	"wikipdf/lib/tasks"
	"net/http"
)

type InsertResult struct {
	JobID string `json:jobID`
}

func (service *Service) InsertHandler(
	writer http.ResponseWriter, request *http.Request) {
	encoder := json.NewEncoder(writer)
	jobID := ksuid.New().String()

	serviceUrl, err := service.run.GetServiceUrl(config.Region(), "wikipdf")
	if err != nil {
		handleError(
			writer,
			"Error getting service URL",
			err,
		)
		return
	}

	task := &tasks.Task{
		Queue:          "pdfJob",
		TaskID:         jobID,
		URL:            fmt.Sprintf("%s/pdfs/%s", serviceUrl, jobID),
		Body:           nil,
		ServiceAccount: config.ServiceAccount(),
	}

	service.WritePending(jobID, task)
	err = service.tasks.AddTask(request.Context(), task)

	if err != nil {
		handleError(
			writer,
			"Error sending task",
			err,
		)
		return
	}

	writer.WriteHeader(http.StatusOK)
	err = encoder.Encode(InsertResult{JobID: jobID})
	if err != nil {
		handleError(
			writer,
			"Error encoding result",
			err,
		)
		return
	}

}
