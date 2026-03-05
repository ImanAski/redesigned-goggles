package transcoder

import (
	"context"
	"io"
	"os/exec"
)

type SteramOpts struct {
	Format  string `json:"format"`
	Bitrate string `json:"bitrate"`
}

func Steram(ctx context.Context, inputPath string, opts SteramOpts, w io.Writer) error {
	if opts.Format == "" {
		opts.Format = "mp3"
	}

	if opts.Bitrate == "" {
		opts.Format = "192k"
	}

	args := []string{
		"-hide_banner", "-loglevel", "error", // Keep logs clean
		"-i", inputPath,
		"-c:a", "libmp3lame",
		"-b:a", opts.Bitrate,
		"-f", opts.Format,
		"pipe:1",
	}

	if inputPath == "anullsrc" {
		args = append(args, "-f", "lavfi")
	}

	// 3. Create Command with Context
	// CommandContext will automatically kill the ffmpeg process if ctx is canceled
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	// 4. Link FFmpeg's Stdout directly to our io.Writer
	cmd.Stdout = w

	// 5. Run it
	// Run() starts the process and waits for it to complete
	return cmd.Run()
}
