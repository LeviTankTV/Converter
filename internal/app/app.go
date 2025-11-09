package app

import (
	"context"
	"video-converter/internal/job"
	"video-converter/internal/service"
)

type App struct {
	converter service.Converter
}

func NewApp(converter service.Converter) *App {
	return &App{
		converter: converter,
	}
}

func (a *App) Run(ctx context.Context, fileID string, format string) (string, error) {
	return a.converter.ProcessJob(ctx, fileID, job.Format(format))
}
