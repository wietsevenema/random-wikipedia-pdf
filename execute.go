package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func (service *Service) ExecuteHandler(
	writer http.ResponseWriter, request *http.Request) {

	//FIXME extract queuename
	queueName := request.Header.Get("X-Cloudtasks-Queuename")
	if queueName != "pdfJob" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(request)
	jobID := vars["jobID"]

	sw := service.bucket.
		Object(jobID + ".pdf").
		NewWriter(context.Background())
	defer sw.Close()

	var buf []byte
	err := printPDF(context.Background(), &buf)
	if err != nil {
		handleError(writer, "Error printing pdf", err)
		return
	}
	_, err = sw.Write(buf)
	if err != nil {
		handleError(writer, "Error writing data", err)
		return
	}

	err = service.bucket.
		Object(jobID + ".pending").Delete(context.Background())
	if err != nil {
		log.Printf("ERROR: %s", err)
		// Not quitting on this error
	}

	writer.WriteHeader(http.StatusOK)
}
