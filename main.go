package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"

	"wikipdf/lib/run"
	"wikipdf/lib/tasks"
	"wikipdf/storage"
)

type Service struct {
	tasks   *tasks.Tasks
	run     *run.Client
	storage *storage.Storage
}

func main() {
	log.Println("Creating Cloud Storage Client")
	storage, err := storage.NewClient()
	if err != nil {
		panic(err)
	}

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
		storage: storage,
		tasks:   tasks,
		run:     runClient}
	service.listenAndServe()
}

func (service *Service) listenAndServe() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "<a href='/pdf'>Random wikipedia PDF</a>")
	})

	jobIDFormat := "{jobID:[a-zA-Z0-9]{22}}"
	r.HandleFunc("/jobs", service.InsertHandler).
		Methods("POST")

	r.HandleFunc("/jobs/"+jobIDFormat, service.StatusHandler).
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

//FIXME update handlers ("appError" approach)
func handleError(writer http.ResponseWriter, msg string, err error) {
	log.Printf("ERROR: %s", err)
	http.Error(
		writer,
		msg,
		http.StatusInternalServerError,
	)
}
