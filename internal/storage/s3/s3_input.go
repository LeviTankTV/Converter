package s3

import (
	"context"
	"errors"
	"io"
	"video-converter/internal/storage"
)

var ErrS3NotImplemented = errors.New("S3 storage not implemented yet. Waiting for infra config")

type S3InputStorage struct{}

func NewS3InputStorage() storage.InputStorage {
	return &S3InputStorage{}
}

func (s *S3InputStorage) Get(_ context.Context, _ string) (io.ReadCloser, error) {
	return nil, ErrS3NotImplemented
}