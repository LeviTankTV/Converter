package storage

import (
	"context"
	"errors"
	"io"
)

var ErrS3NotImplemented = errors.New("S3 storage not implemented yet. Waiting for infra config")

type S3InputStorage struct{}

func (s *S3InputStorage) Get(_ context.Context, _ string) (io.ReadCloser, error) {
	return nil, ErrS3NotImplemented
}