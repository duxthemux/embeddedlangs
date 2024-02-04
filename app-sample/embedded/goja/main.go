package main

import (
	"context"
	"log"

	"github.com/dop251/goja"
)

func Some(in map[string]any) {
	in["FromGoja"] = "Changed by Goja"
	log.Printf("Called from JS")
}

func run(ctx context.Context) error {
	vm := goja.New()
	in := map[string]any{}
	vm.Set("pin", in)
	vm.Set("Some", Some)
	v, err := vm.RunString(`
pin.INJS="In JS";
Some(pin);
pin;
`)
	if err != nil {
		panic(err)
	}
	num := v.Export()

	log.Printf("%#v", num)

	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		panic(err)
	}
}
