package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Engine struct {
	Frequency   time.Duration
	Repo        string
	LastHash    string
	lastRepoDir string
}

func (e *Engine) Run(ctx context.Context) error {
	if e.Frequency == 0 {
		e.Frequency = time.Minute
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	slog.Info("running from", "dir", wd)

	parts := strings.Split(e.Repo, "/")
	e.lastRepoDir = parts[len(parts)-1]
	e.lastRepoDir = strings.Split(e.lastRepoDir, ".")[0]
	e.lastRepoDir, err = filepath.Abs(filepath.Join(wd, e.lastRepoDir))
	if err != nil {
		return err
	}

	_, err = os.Stat(e.lastRepoDir)
	if errors.Is(err, os.ErrNotExist) {
		if err = e.Clone(ctx); err != nil {
			return err
		}
		e.CheckHash(ctx)
	}

	for {
		err = e.Pull(ctx)
		if err != nil {
			return err
		}
		if err = e.CheckHash(ctx); err != nil {
			return err
		}
		time.Sleep(e.Frequency)
	}
}

func (e *Engine) Clone(ctx context.Context) error {
	bs, err := exec.Command("sh", "-c", "git clone "+e.Repo).CombinedOutput()
	os.Stdout.Write(bs)
	if bytes.Contains(bs, []byte("already exists and is not an empty directory")) {
		return nil
	}

	return err
}

func (e *Engine) Pull(ctx context.Context) error {
	bs, err := exec.Command("sh", "-c", fmt.Sprintf("git -C %s remote add origin %s", e.lastRepoDir, e.Repo)).CombinedOutput()
	os.Stdout.Write(bs)

	bs, err = exec.Command("sh", "-c", fmt.Sprintf("git -C %s pull", e.lastRepoDir)).CombinedOutput()
	os.Stdout.Write(bs)

	return err
}

func (e *Engine) CheckHash(ctx context.Context) error {
	bs, err := exec.Command("sh", "-c", fmt.Sprintf("git -C %s rev-parse HEAD", e.lastRepoDir)).CombinedOutput()
	if err != nil {
		return err
	}

	lh := strings.TrimSpace(string(bs))
	if lh != e.LastHash {
		e.LastHash = lh
		slog.Info("Got new hash", "hash", lh)

	}

	return nil
}
