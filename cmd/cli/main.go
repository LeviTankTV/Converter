package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"
	"video-converter/internal/app"
	"video-converter/internal/cli"
	"video-converter/internal/service"
	"video-converter/internal/storage/local"
)

func main() {
	// Парсим аргументы
	flag.Parse()
	args := cli.Parse()

	if args.FileID == "" || args.Format == "" {
		flag.Usage()
		log.Fatal("Требуются --file_id и --format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	// Вайринг зависимостей
	inputStorage := local.NewLocalInputStorage("./uploads")
	outputStorage := local.NewLocalOutputStorage("./outputs")
	converterService := service.NewVideoService(inputStorage, outputStorage)
	application := app.NewApp(converterService)

	result, err := application.Run(ctx, args.FileID, args.Format)
	if err != nil {
		log.Fatalf("❌ Конвертация упала: %v", err)
		os.Exit(1)
	}

	log.Printf("✅ Готово! Результат: %s", result)
}
