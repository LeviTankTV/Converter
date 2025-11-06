package converter

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"video-converter/internal/job"
)

func Convert(ctx context.Context, input io.Reader, output io.Writer, format job.Format) error {
	args := buildFFmpegArgs(format)
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	cmd.Stdin = input
	cmd.Stdout = output
	// cmd.Stderr = os.Stderr // раскомментируй для отладки

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}
	return nil
}

func buildFFmpegArgs(format job.Format) []string {
	container := formatToContainer(format)
	base := []string{"-i", "pipe:0", "-y", "-f", container, "pipe:1"}

	switch format {
	case job.MP4_H264:
		return append(base,
			"-c:v", "libx264",
			"-preset", "fast",
			"-crf", "23",
			"-c:a", "aac",
			"-b:a", "128k",
		)
	case job.MP4_H265:
		return append(base,
			"-c:v", "libx265",
			"-crf", "28",
			"-c:a", "aac",
			"-b:a", "128k",
		)
	case job.WEBM:
		return append(base,
			"-c:v", "libvpx-vp9",
			"-crf", "30",
			"-b:v", "0",
			"-c:a", "libopus",
			"-b:a", "128k",
		)
	case job.GIF:
		// Для GIF нужна палитра → требует файл. Обход: временно сохраняем.
		// В этом примере GIF обрабатывается через файл (см. CLI).
		// Альтернатива — использовать filter_complex с палитрой в памяти (сложно).
		panic("GIF requires file-based processing — handle in CLI")
	case job.MP3:
		return append(base,
			"-vn",
			"-c:a", "libmp3lame",
			"-b:a", "192k",
		)
	case job.AAC:
		return append(base,
			"-vn",
			"-c:a", "aac",
			"-b:a", "128k",
		)
	case job.OGG:
		return append(base,
			"-vn",
			"-c:a", "libopus",
			"-b:a", "128k",
		)
	default:
		panic("unsupported format: " + string(format))
	}
}

func formatToContainer(f job.Format) string {
	switch f {
	case job.MP4_H264, job.MP4_H265, job.AAC:
		return "mp4"
	case job.WEBM:
		return "webm"
	case job.GIF:
		return "gif"
	case job.MP3:
		return "mp3"
	case job.OGG:
		return "ogg"
	default:
		return "mp4"
	}
}
