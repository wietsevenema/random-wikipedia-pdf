package main

import (
	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func (service *Service) ResultHandler(
	writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	jobID := vars["jobID"]

	reader, err := service.storage.ResultReader(
		request.Context(),
		jobID,
	)
	if err == storage.ErrObjectNotExist {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		handleError(writer, "Error reading file", err)
		return
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}
}
