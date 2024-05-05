package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/duxthemux/embeddedlangs/app-sample/git-clone/internal"
)

func run(ctx context.Context) error {
	repo := os.Getenv("GITHUB_REPOSITORY")
	freqStr := os.Getenv("GIT_PUSH_FREQUENCY")
	wd := os.Getenv("GIT_CLONE_DIR")
	if wd == "" {
		wd = "."
	}
	freq, err := time.ParseDuration(freqStr)
	if err != nil {
		slog.Warn("Cant part of GIT_PUSH_FREQUENCY - using default 1m", "err", err)
		freq = 1 * time.Minute
	}

	err = os.Chdir(wd)
	if err != nil {
		return err
	}

	gce := internal.Engine{
		Repo:      repo,
		Frequency: freq,
	}

	return gce.Run(ctx)
}
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}
