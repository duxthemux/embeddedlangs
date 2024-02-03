package main

import (
	"context"
	"fmt"
	"log"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/token"
)

type MyGoFunc struct {
}

func (m *MyGoFunc) String() string {
	return ""
}

func (m *MyGoFunc) BinaryOp(op token.Token, rhs tengo.Object) (tengo.Object, error) {
	return nil, nil
}

func (m *MyGoFunc) IsFalsy() bool {
	return false
}

func (m *MyGoFunc) Equals(another tengo.Object) bool {
	return false
}

func (m *MyGoFunc) Copy() tengo.Object {
	return m
}

func (m *MyGoFunc) IndexGet(index tengo.Object) (value tengo.Object, err error) {
	return nil, nil
}

func (m *MyGoFunc) IndexSet(index, value tengo.Object) error {
	return nil
}

func (m *MyGoFunc) Iterate() tengo.Iterator {
	return nil
}

func (m *MyGoFunc) CanIterate() bool {
	return false
}

func (m *MyGoFunc) TypeName() string {
	return "MyGoFunc"
}

func (m *MyGoFunc) CanCall() bool {
	return true
}

func (m *MyGoFunc) Call(args ...tengo.Object) (ret tengo.Object, err error) {
	log.Printf("Hello from called: %#v", args)

	return args[0], nil
}

func run(ctx context.Context) error {
	// create a new Script instance
	script := tengo.NewScript([]byte(
		`each := func(seq, fn) {
    for x in seq { fn(x) }
}

ret := fnFromGo(a)

sum := 0
mul := 1
each([a, b, c, d], func(x) {
    sum += x
    mul *= x
})`))

	// set values
	_ = script.Add("a", 1)
	_ = script.Add("b", 9)
	_ = script.Add("c", 8)
	_ = script.Add("d", 4)
	err := script.Add("fnFromGo", &MyGoFunc{})
	if err != nil {
		return err
	}

	// run the script
	compiled, err := script.RunContext(context.Background())
	if err != nil {
		panic(err)
	}

	// retrieve values
	sum := compiled.Get("sum")
	mul := compiled.Get("mul")
	ret := compiled.Get("ret")
	fmt.Println(sum, mul, ret) // "22 288"
	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		panic(err)
	}
}
