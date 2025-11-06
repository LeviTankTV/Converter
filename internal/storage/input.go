package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type InputStorage interface {
	Get(ctx context.Context, fileID string) (io.ReadCloser, error)
}

type LocalInputStorage struct {
	BaseDir string
}

func (s *LocalInputStorage) Get(_ context.Context, fileID string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(s.BaseDir, fileID))
}
