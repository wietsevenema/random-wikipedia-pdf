package main

import (
	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func (service *Service) ResultHandler(
	writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	jobID := vars["jobID"]

	reader, err := service.bucket.Object(jobID + ".pdf").NewReader(request.Context());
	if err == nil {
		_, err := io.Copy(writer, reader)
		if err != nil {
			handleError(writer, "Error writing response", err)
			return
		}
		writer.WriteHeader(http.StatusOK)
		return
	}

	if err == storage.ErrObjectNotExist {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	handleError(writer, "Error reading file", err)
}
