package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ffandroid/autoprime/services"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <tmdb_id> <season> <episode>\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}

	tmdbID, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil || tmdbID <= 0 {
		fmt.Fprintln(os.Stderr, "tmdb_id must be a positive integer")
		os.Exit(2)
	}

	season, err := strconv.Atoi(os.Args[2])
	if err != nil || season < 0 {
		fmt.Fprintln(os.Stderr, "season must be a non-negative integer")
		os.Exit(2)
	}

	episode, err := strconv.Atoi(os.Args[3])
	if err != nil || episode < 0 {
		fmt.Fprintln(os.Stderr, "episode must be a positive integer")
		os.Exit(2)
	}

	services.AddVideo(tmdbID, season, episode)
}
