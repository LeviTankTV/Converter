package local

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"video-converter/internal/job"
	"video-converter/internal/storage"
)

type LocalOutputStorage struct {
	BaseDir string
}

func NewLocalOutputStorage(baseDir string) storage.OutputStorage {
	return &LocalOutputStorage{BaseDir: baseDir}
}

func (s *LocalOutputStorage) Put(_ context.Context, jobID string, format job.Format) (io.WriteCloser, string, error) {
	ext := map[job.Format]string{
		job.MP4_H264: ".mp4",
		job.MP4_H265: ".mp4",
		job.WEBM:     ".webm",
		job.GIF:      ".gif",
		job.MP3:      ".mp3",
		job.AAC:      ".m4a",
		job.OGG:      ".ogg",
	}[format]

	filename := jobID + ext
	path := filepath.Join(s.BaseDir, filename)
	file, err := os.Create(path)
	return file, filename, err
}