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

	err = os.WriteFile("run.go", []byte(`package main

import (
	"encoding/json"
	"os"
)

func main() {
	m := map[string]any{}
	if err := json.NewDecoder(os.Stdin).Decode(&m); err != nil {
		panic(err)
	}
	m["fromGo"] = "Came from Go!!!"
	if err := json.NewEncoder(os.Stdout).Encode(m); err != nil {
		panic(err)
	}
}`),

		0o600)
	if err != nil {
		return err
	}

	defer os.Remove("run.go")

	cmd := exec.Command("go", "run", "run.go")

	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	serr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go io.Copy(os.Stderr, serr)

	go io.Copy(os.Stdout, out)

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
