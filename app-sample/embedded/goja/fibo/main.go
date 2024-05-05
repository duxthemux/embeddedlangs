package main

import (
	"context"
	_ "embed"
	"log"
	"time"

	"github.com/dop251/goja"
)

//go:embed fibt.js
var src string

func main() {
	start := time.Now()
	if err := run(context.Background()); err != nil {
		panic(err)
	}
	dur := time.Since(start)
	log.Printf("took %v", dur)
}

func Some(in map[string]any) {
	in["FromGoja"] = "Changed by Goja"
	log.Printf("Called from JS")
}

func run(ctx context.Context) error {
	vm := goja.New()
	v, err := vm.RunString(src)
	if err != nil {
		panic(err)
	}
	num := v.Export()

	log.Printf("%#v", num)

	return nil
}
