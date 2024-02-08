package main

import (
	"context"
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func run(ctx context.Context) error {
	thread := &starlark.Thread{Name: "my thread", Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) }}
	_, prog, err := starlark.SourceProgramOptions(&syntax.FileOptions{}, "main.star", `
def fibonacci(n):
    res = list(range(n))
    for i in res[2:]:
        res[i] = res[i-2] + res[i-1]
    return res

print(fibonacci(10))
`, func(s string) bool {
		return false
	})

	if err != nil {
		return err
	}

	_, err = prog.Init(thread, starlark.StringDict{})

	return err
}

func main() {
	if err := run(context.Background()); err != nil {
		panic(err)
	}
}
