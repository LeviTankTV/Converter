// internal/service/service.go
package service

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"video-converter/internal/converter"
	"video-converter/internal/job"
	"video-converter/internal/storage"
)

type VideoService struct {
	inputStorage  storage.InputStorage
	outputStorage storage.OutputStorage
}

func NewVideoService(input storage.InputStorage, output storage.OutputStorage) *VideoService {
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
		// Очистка при ошибке (для локального режима)
		os.Remove(filepath.Join("outputs", resultName))
		return "", err
	}

	return resultName, nil
}

func (s *VideoService) processGIF(ctx context.Context, jobID string, fileID string) (string, error) {
	// 1. Скачиваем исходник во временный файл
	tmpInput, err := ioutil.TempFile("", "input-*.tmp")
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

	// 2. Генерация палитры
	paletteFile, err := ioutil.TempFile("", "palette-*.png")
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
	tmpOutput, err := ioutil.TempFile("", "output-*.gif")
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
		os.Remove(filepath.Join("outputs", resultName))
		return "", err
	}

	return resultName, nil
}
