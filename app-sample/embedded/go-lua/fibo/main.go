package main

import (
	"context"
	_ "embed"
	"log"
	"time"

	lua "github.com/yuin/gopher-lua"
)

//go:embed fibt.lua
var src string

func run(ctx context.Context) error {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(src); err != nil {
		panic(err)
	}
	return nil
}

func main() {
	start := time.Now()
	if err := run(context.Background()); err != nil {
		panic(err)
	}
	dur := time.Since(start)
	log.Printf("took %v", dur)
}
