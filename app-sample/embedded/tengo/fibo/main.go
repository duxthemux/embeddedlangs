package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"time"

	"github.com/d5/tengo/v2"
)

//go:embed fibt.tengo
var src []byte

func run(ctx context.Context) error {
	script := tengo.NewScript(src)

	// run the script
	compiled, err := script.RunContext(ctx)
	if err != nil {
		panic(err)
	}

	ret := compiled.Get("ret")
	fmt.Printf("%v", ret) //
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
