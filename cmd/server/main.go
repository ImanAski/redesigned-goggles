package main

import (
	"encoding/json"
	"go-music/internal/transcoder"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		opts := transcoder.SteramOpts{Format: "mp3", Bitrate: "128k"}

		// Set headers so the browser knows it's an audio stream
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Transfer-Encoding", "chunked")

		err := transcoder.Steram(r.Context(), "./internal/transcoder/test.mp3", opts, w)
		if err != nil {
			log.Println(err)
			return
		}
	})

	http.HandleFunc("/stream/metadata", func(w http.ResponseWriter, r *http.Request) {
		meta, err := transcoder.GetMetadata(r.Context(), "./internal/transcoder/test.mp3")
		if err != nil {
			http.Error(w, "Failed to extract metadata", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(meta)
	})

	log.Fatalln(http.ListenAndServe(":8000", nil))
}
