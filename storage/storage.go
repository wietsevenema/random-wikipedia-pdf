package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"log"
	"wikipdf/lib/config"
	"wikipdf/lib/tasks"
)

// Tasks is the client object that contains
// pointers to the backend clients
type Storage struct {
	bucket *storage.BucketHandle
}

const inProgressSuffix = ".pending"
const pdfSuffix = ".pdf"

func NewClient() (*Storage, error) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create Cloud Storage client: %v", err)
	}
	bucket := client.Bucket(config.ProjectID() + "-pdf")

	if err != nil {
		return nil, err
	}
	return &Storage{bucket: bucket}, nil
}

func (s *Storage) WriteInProgress(
	jobID string,
	task *tasks.Task,
) error {
	sw := s.bucket.
		Object(jobID + inProgressSuffix).
		NewWriter(context.Background())
	defer sw.Close()

	en := json.NewEncoder(sw)
	return en.Encode(task)
}

func (s *Storage) Completed(ctx context.Context, jobID string) (bool, error) {
	_, err := s.bucket.Object(jobID + pdfSuffix).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Storage) InProgress(
	ctx context.Context,
	jobID string,
) (bool, error) {

	_, err := s.bucket.Object(
		jobID + inProgressSuffix,
	).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Storage) RemoveInProgress(jobID string) error {
	return s.bucket.
		Object(jobID + inProgressSuffix).
		Delete(context.Background())

}

func (s *Storage) ResultWriter(
	ctx context.Context,
	jobID string,
) *storage.Writer {

	return s.bucket.
		Object(jobID + ".pdf").
		NewWriter(context.Background())
}

func (s *Storage) ResultReader(
	ctx context.Context,
	jobID string,
) (*storage.Reader, error) {
	return s.bucket.Object(jobID + pdfSuffix).NewReader(ctx)
}
