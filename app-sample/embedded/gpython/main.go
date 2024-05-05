package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/go-python/gpython/py"
)

func init() {

	// For each of your embedded python types, attach instance methods.
	// When an instance method is invoked, the "self" py.Object is the instance.
	PyVacationStopType.Dict["Set"] = py.MustNewMethod("Set", VacationStop_Set, 0, "")
	PyVacationStopType.Dict["Get"] = py.MustNewMethod("Get", VacationStop_Get, 0, "")
	PyVacationType.Dict["add_stops"] = py.MustNewMethod("Vacation.add_stops", Vacation_add_stops, 0, "")
	PyVacationType.Dict["num_stops"] = py.MustNewMethod("Vacation.num_stops", Vacation_num_stops, 0, "")
	PyVacationType.Dict["get_stop"] = py.MustNewMethod("Vacation.get_stop", Vacation_get_stop, 0, "")

	// Bind methods attached at the module (global) level.
	// When these are invoked, the first py.Object param (typically "self") is the bound *Module instance.
	methods := []*py.Method{
		py.MustNewMethod("VacationStop_new", VacationStop_new, 0, ""),
		py.MustNewMethod("Vacation_new", Vacation_new, 0, ""),
	}

	// Register a ModuleImpl instance used by the gpython runtime to instantiate new py.Module when first imported.
	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "mylib_go",
			Doc:  "Example embedded python module",
		},
		Methods: methods,
		Globals: py.StringDict{
			"PY_VERSION": py.String("Python 3.4 (github.com/go-python/gpython)"),
			"GO_VERSION": py.String(fmt.Sprintf("%s on %s %s", runtime.Version(), runtime.GOOS, runtime.GOARCH)),
			"MYLIB_VERS": py.String("Vacation 1.0 by Fletch F. Fletcher"),
		},
		OnContextClosed: func(instance *py.Module) {
			py.Println(instance, "<<< host py.Context of py.Module instance closing >>>\n+++")
		},
	})
}

var (
	PyVacationStopType = py.NewType("Stop", "")
	PyVacationType     = py.NewType("Vacation", "")
)

func run(ctx context.Context) error {

	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		panic(err)
	}
}
