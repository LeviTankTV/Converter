package storage

import (
	"context"
	"io"
)

type S3OutputStorage struct{}

func (s *S3OutputStorage) Put(_ context.Context, _, _ string) (io.WriteCloser, string, error) {
	return nil, "", ErrS3NotImplemented
}