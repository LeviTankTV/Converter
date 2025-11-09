package service

import (
	"context"
	"io"
	"os"
	"os/exec"
	"video-converter/internal/converter"
	"video-converter/internal/job"
	"video-converter/internal/storage"
)

type Converter interface {
	ProcessJob(ctx context.Context, fileID string, format job.Format) (string, error)
}

type VideoService struct {
	inputStorage  storage.InputStorage
	outputStorage storage.OutputStorage
}

func NewVideoService(input storage.InputStorage, output storage.OutputStorage) Converter {
	return &VideoService{
		inputStorage:  input,
		outputStorage: output,
	}
}

func (s *VideoService) ProcessJob(ctx context.Context, fileID string, format job.Format) (string, error) {
	jobID := "job-" + fileID

	if format == job.GIF {
		return s.processGIF(ctx, jobID, fileID)
	}

	input, err := s.inputStorage.Get(ctx, fileID)
	if err != nil {
		return "", err
	}
	defer input.Close()

	output, resultName, err := s.outputStorage.Put(ctx, jobID, format)
	if err != nil {
		return "", err
	}
	defer output.Close()

	if err := converter.Convert(ctx, input, output, format); err != nil {
		return "", err
	}

	return resultName, nil
}

func (s *VideoService) processGIF(ctx context.Context, jobID string, fileID string) (string, error) {
	// 1. Скачиваем исходник во временный файл
	tmpInput, err := os.CreateTemp("", "input-*.tmp")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpInput.Name())
	defer tmpInput.Close()

	src, err := s.inputStorage.Get(ctx, fileID)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(tmpInput, src); err != nil {
		src.Close()
		return "", err
	}
	src.Close()
	tmpInput.Seek(0, io.SeekStart)

	// 2. Генерация палитры
	paletteFile, err := os.CreateTemp("", "palette-*.png")
	if err != nil {
		return "", err
	}
	defer os.Remove(paletteFile.Name())
	palettePath := paletteFile.Name()
	paletteFile.Close()

	palCmd := exec.CommandContext(ctx, "ffmpeg", "-i", tmpInput.Name(), "-vf", "palettegen", "-y", palettePath)
	if err := palCmd.Run(); err != nil {
		return "", err
	}

	// 3. Конвертация с палитрой
	tmpOutput, err := os.CreateTemp("", "output-*.gif")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpOutput.Name())
	outputPath := tmpOutput.Name()
	tmpOutput.Close()

	convCmd := exec.CommandContext(ctx, "ffmpeg", "-i", tmpInput.Name(), "-i", palettePath, "-lavfi", "paletteuse", "-y", outputPath)
	if err := convCmd.Run(); err != nil {
		return "", err
	}

	// 4. Загружаем результат в OutputStorage
	output, resultName, err := s.outputStorage.Put(ctx, jobID, job.GIF)
	if err != nil {
		return "", err
	}
	defer output.Close()

	tmpResult, err := os.Open(outputPath)
	if err != nil {
		return "", err
	}
	defer tmpResult.Close()

	if _, err := io.Copy(output, tmpResult); err != nil {
		return "", err
	}

	return resultName, nil
}
