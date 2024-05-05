package main

import (
	"context"
	_ "embed"
	"log"
	"time"

	"github.com/go-python/gpython/py"
	_ "github.com/go-python/gpython/stdlib"
)

//go:embed fibt.py
var src string

func main() {
	start := time.Now()
	if err := run(context.Background()); err != nil {
		panic(err)
	}
	dur := time.Since(start)
	log.Printf("took %v", dur)
}

func run(goCtx context.Context) error {
	// See type Context interface and related docs
	ctx := py.NewContext(py.DefaultContextOpts())

	// This drives modules being able to perform cleanup and release resources
	defer ctx.Close()

	_, err := py.RunSrc(ctx, src, "", nil)
	// _, err := py.RunFile(ctx, "fibt.py", py.CompileOpts{}, nil)
	if err != nil {
		return err
	}

	return nil
}
