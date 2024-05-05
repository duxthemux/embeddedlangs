package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"time"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

//go:embed fibt.star
var src string

func main() {
	start := time.Now()
	if err := run(context.Background()); err != nil {
		panic(err)
	}
	dur := time.Since(start)
	log.Printf("took %v", dur)
}

func Fufu(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.String("fufufufu"), nil
}

func run(ctx context.Context) error {
	builtins := starlark.StringDict{
		"fufu": starlark.NewBuiltin("fufu", Fufu),
	}

	thread := &starlark.Thread{
		Name:  "my thread",
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
	}

	_, prog, err := starlark.SourceProgramOptions(&syntax.FileOptions{
		Recursion: true,
	}, "main.star", src, func(s string) bool {
		if s == "fufu" {
			return true
		}
		return false
	})

	if err != nil {
		return err
	}

	_, err = prog.Init(thread, builtins)

	return err
}
