package main

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/gorilla/mux"
	"wikipdf/lib/config"
	"wikipdf/lib/run"
	"wikipdf/lib/tasks"
	"log"
	"net/http"
	"os"
)

type Service struct {
	bucket *storage.BucketHandle
	tasks  *tasks.Tasks
	run    *run.Client
}

func main() {

	log.Println("Creating Cloud Storage Client")
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Cloud Storage client: %v", err)
	}
	bucket := client.Bucket(config.ProjectID() + "-pdf")

	// Tasks
	log.Println("Creating Cloud Tasks Client")
	tasks, err := tasks.NewClient()
	if err != nil {
		panic(err)
	}

	// Cloud Run Service Map
	log.Println("Creating Cloud Run Service Map")
	runClient, err := run.NewClient()
	if err != nil {
		panic(err)
	}

	service := Service{
		bucket: bucket,
		tasks:  tasks,
		run:    runClient}
	service.listenAndServe()
}

func (service *Service) listenAndServe() {
	r := mux.NewRouter()

	jobIDFormat := "{jobID:[a-zA-Z0-9]{27}}"
	r.HandleFunc("/jobs", service.InsertHandler).
		Methods("POST")

	r.HandleFunc("/jobs/"+jobIDFormat, service.StateHandler).
		Methods("GET")

	r.HandleFunc("/pdfs/"+jobIDFormat, service.ResultHandler).
		Methods("GET")

	r.HandleFunc("/pdfs/"+jobIDFormat, service.ExecuteHandler).
		Methods("POST")

	r.HandleFunc("/pdf", service.PdfHandler).
		Methods("GET")

	// Start the HTTP server on PORT
	// Grab the PORT from the environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func handleError(writer http.ResponseWriter, msg string, err error) {
	log.Printf("ERROR: %s", err)
	http.Error(
		writer,
		msg,
		http.StatusInternalServerError,
	)
}
