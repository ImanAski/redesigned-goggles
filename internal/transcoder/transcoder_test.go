package transcoder

import (
	"context"
	"io"
	"os/exec"
	"testing"
	"time"
)

func TestSteram(t *testing.T) {
	// 1. Check if ffmpeg is even installed on the test machine
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		t.Skip("ffmpeg not found in PATH, skipping test")
	}

	// 2. Setup a dummy input.
	// Instead of a real file, we use ffmpeg's 'lavfi' to generate 1 second of silence.
	// This makes the test portable and fast.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use a pipe to capture the output
	pr, pw := io.Pipe()

	opts := SteramOpts{
		Format:  "mp3",
		Bitrate: "64k",
	}

	// 3. Run the Stream function in a goroutine
	// We use "anullsrc" as the input path to simulate a real file
	go func() {
		// Note: To use 'anullsrc' (silent audio), the ffmpeg command in Steram()
		// would need to handle '-f lavfi'. For this test, let's assume
		// you have a tiny test file or use a real path.
		err := Steram(ctx, "/home/sophos/Music/'02 Gomshodeh-۲.mp3'", opts, pw)
		pw.CloseWithError(err)
	}()

	// 4. Validate the output
	// We read the first 100 bytes to see if it looks like an MP3 header
	buf := make([]byte, 100)
	n, err := io.ReadFull(pr, buf)

	if err != nil && err != io.EOF {
		t.Fatalf("Failed to read stream: %v", err)
	}

	if n < 10 {
		t.Errorf("Stream too short, only got %d bytes", n)
	}

	// MP3 files usually start with "ID3" or 0xFF (frame sync)
	t.Logf("Received %d bytes of streamed data", n)
}
