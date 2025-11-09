package s3

import (
	"context"
	"io"
	"video-converter/internal/job"
	"video-converter/internal/storage"
)

type S3OutputStorage struct{}

func NewS3OutputStorage() storage.OutputStorage {
	return &S3OutputStorage{}
}

func (s *S3OutputStorage) Put(_ context.Context, jobID string, format job.Format) (io.WriteCloser, string, error) {
	return nil, "", ErrS3NotImplemented
}
