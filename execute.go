package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"wikipdf/lib/config"
)

func (service *Service) ExecuteHandler(
	writer http.ResponseWriter, request *http.Request) {

	queueName := request.Header.Get("X-Cloudtasks-Queuename")
	if queueName != config.QueueName() {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(request)
	jobID := vars["jobID"]

	inProgress, err := service.storage.InProgress(
		request.Context(),
		jobID,
	)
	if err != nil {
		handleError(writer, "Error reading job state", err)
		return
	}
	if !inProgress {
		http.NotFound(writer, request)
		return
	}

	sw := service.storage.ResultWriter(request.Context(), jobID)
	defer sw.Close()

	var buf []byte
	err = printPDF(context.Background(), &buf)
	if err != nil {
		handleError(writer, "Error printing pdf", err)
		return
	}
	_, err = sw.Write(buf)
	if err != nil {
		handleError(writer, "Error writing data", err)
		return
	}

	err = service.storage.RemoveInProgress(jobID)
	if err != nil {
		log.Printf("ERROR: %s", err)
		// Not quitting on this error
	}

	writer.WriteHeader(http.StatusOK)
}
