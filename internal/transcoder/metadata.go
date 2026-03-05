// internal/transcoder/metadata.go
package transcoder

import (
	"context"
	"encoding/json"
	"os/exec"
	"strconv"
)

type Metadata struct {
	Title    string  `json:"title"`
	Artist   string  `json:"artist"`
	Duration float64 `json:"duration"`
	Format   string  `json:"format"`
}

func GetMetadata(ctx context.Context, inputPath string) (*Metadata, error) {
	// ffprobe -v quiet -print_format json -show_format input.mp3
	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		inputPath,
	}

	out, err := exec.CommandContext(ctx, "ffprobe", args...).Output()
	if err != nil {
		return nil, err
	}

	// Internal struct to match ffprobe's deep JSON output
	var probe struct {
		Format struct {
			Duration string `json:"duration"`
			Tags     struct {
				Title  string `json:"title"`
				Artist string `json:"artist"`
			} `json:"tags"`
			FormatName string `json:"format_name"`
		} `json:"format"`
	}

	if err := json.Unmarshal(out, &probe); err != nil {
		return nil, err
	}

	duration, _ := strconv.ParseFloat(probe.Format.Duration, 64)

	return &Metadata{
		Title:    probe.Format.Tags.Title,
		Artist:   probe.Format.Tags.Artist,
		Duration: duration,
		Format:   probe.Format.FormatName,
	}, nil
}
