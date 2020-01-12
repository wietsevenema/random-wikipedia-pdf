package tasks

import (
	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta3"
	"context"
	"fmt"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta3"
	"log"
	"wikipdf/lib/config"
)

// Tasks is the client object that contains
// pointers to the backend clients
type Tasks struct {
	tasks *cloudtasks.Client
}

func NewClient() (*Tasks, error) {
	tasks, err := cloudtasks.NewClient(context.Background())
	if err != nil {
		return nil, err
	}
	return &Tasks{tasks}, nil
}

type Task struct {
	Queue          string `json:queue`
	TaskID         string `json:taskID`
	URL            string `json:url`
	Body           []byte `json:"-"`
	ServiceAccount string `json:serviceAccount`
}

func (client *Tasks) AddTask(
	ctx context.Context,
	task *Task) error {

	log.Printf("submitting task to %s", task.URL)

	// Build the HTTP request recipe
	requestRecipe := &taskspb.HttpRequest{
		// The URL to send the request to
		Url:        task.URL,
		Body:       task.Body,
		HttpMethod: taskspb.HttpMethod_POST,
		AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
			// This tells Cloud Tasks to add an Authorization
			// header with the identity of the calling function
			OidcToken: &taskspb.OidcToken{
				ServiceAccountEmail: task.ServiceAccount,
			},
		},
	}

	err := client.sendTask(ctx, task.TaskID, task.Queue, requestRecipe)
	if err != nil {
		return err
	}
	return nil
}

func (client *Tasks) sendTask(
	ctx context.Context,
	ID, queue string,
	requestRecipe *taskspb.HttpRequest) error {

	// Build the queue name
	queueName := client.queueName(queue)

	// Build a unique task name
	taskName := fmt.Sprintf("%v/tasks/%v", queueName, ID)

	// Construct the request to the Cloud Tasks API
	req := &taskspb.CreateTaskRequest{
		Parent: queueName,
		Task: &taskspb.Task{
			Name: taskName,
			PayloadType: &taskspb.Task_HttpRequest{
				HttpRequest: requestRecipe,
			},
		},
	}

	// Send the request to the Cloud Tasks API
	_, err := client.tasks.CreateTask(ctx, req)
	return err
}

func (client *Tasks) queueName(queue string) string {
	queueName := fmt.Sprintf("projects/%v/locations/%v/queues/%v",
		config.ProjectID(),
		config.Region(),
		queue)
	return queueName
}
