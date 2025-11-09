package local

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"video-converter/internal/storage"
)

type LocalInputStorage struct {
	BaseDir string
}

func NewLocalInputStorage(baseDir string) storage.InputStorage {
	return &LocalInputStorage{BaseDir: baseDir}
}

func (s *LocalInputStorage) Get(_ context.Context, fileID string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(s.BaseDir, fileID))
}
