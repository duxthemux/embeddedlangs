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

	err = os.WriteFile("main.py", []byte(`import sys, json;
data = json.load(sys.stdin)

data["field3"]="nono"

print(json.dumps(data))

sys.stdout.write("{}")`), 0o600)
	if err != nil {
		return err
	}

	defer os.Remove("main.py")

	cmd := exec.Command("python3", "main.py")

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
