package main

import (
	"context"
	"log"

	"github.com/yuin/gopher-lua"
)

func DoSome(L *lua.LState) int {
	lv := L.ToString(1)
	log.Printf(lv)
	return 0 /* number of results */
}

func run(ctx context.Context) error {
	L := lua.NewState()
	defer L.Close()
	L.SetGlobal("do_some", L.NewFunction(DoSome))
	if err := L.DoString(`do_some("hello")`); err != nil {
		panic(err)
	}
	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		panic(err)
	}
}
