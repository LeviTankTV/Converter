package storage

import (
	"context"
	"io"
	"video-converter/internal/job"
)

type InputStorage interface {
	Get(ctx context.Context, fileID string) (io.ReadCloser, error)
}

type OutputStorage interface {
	Put(ctx context.Context, jobID string, format job.Format) (io.WriteCloser, string, error)
}