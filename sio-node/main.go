package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"os/exec"
)

type SomeStruct struct {
	A string `json:"a,omitempty"`
	B int    `json:"b,omitempty"`
}

func run(ctx context.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	log.Printf("Running at: %s", wd)
	cmd := exec.Command("node", "main.js")

	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	go func() {
		io.Copy(os.Stdout, out)
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}
	json.NewEncoder(in).Encode(SomeStruct{
		A: "Some fancy msg in",
		B: 12,
	})

	in.Close()

	return cmd.Wait()
}

func main() {
	if err := run(context.Background()); err != nil {
		panic(err)
	}
}
