package cli

import "flag"

type Args struct {
	FileID string
	Format string
}

func Parse() Args {
	var (
		fileID = flag.String("file_id", "", "Имя файла (ключ в хранилище)")
		format = flag.String("format", "", "Формат: mp4_h264, webm, mp3, gif и т.д.")
	)

	return Args{
		FileID: *fileID,
		Format: *format,
	}
}