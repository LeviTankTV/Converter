// cmd/cli/run.go
package cli

import (
	"context"
	"flag"
	"log"
	"time"
	"video-converter/internal/job"
	"video-converter/internal/service"
	"video-converter/internal/storage"
)

func Run() {
	var (
		fileID = flag.String("file_id", "", "Имя файла (ключ в хранилище)")
		format = flag.String("format", "", "Формат: mp4_h264, webm, mp3, gif и т.д.")
	)
	flag.Parse()

	if *fileID == "" || *format == "" {
		flag.Usage()
		log.Fatal("Требуются --file_id и --format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	// БАЗОВАЯ ИНИЦИАЛИЗАЦИЯ (сегодня - локальные файлы, завтра - S3)
	inputStorage := &storage.LocalInputStorage{BaseDir: "./uploads"}
	outputStorage := &storage.LocalOutputStorage{BaseDir: "./outputs"}

	// СОЗДАЁМ СЕРВИС С ИНЪЕКЦИЕЙ ЗАВИСИМОСТЕЙ
	videoService := service.NewVideoService(inputStorage, outputStorage)

	resultName, err := videoService.ProcessJob(ctx, *fileID, job.Format(*format))
	if err != nil {
		log.Fatalf("❌ Конвертация упала: %v", err)
	}

	log.Printf("✅ Готово! Результат: ./outputs/%s", resultName)
}
