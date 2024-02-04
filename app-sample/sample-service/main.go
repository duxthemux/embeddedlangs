package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/dop251/goja"
)

type Out struct {
	Result      any       `json:"result,omitempty"`
	FormulaHash string    `json:"formaula-hash,omitempty"`
	EngineHash  string    `json:"engineHash,omitempty"`
	Input       any       `json:"input,omitempty"`
	User        string    `json:"user,omitempty"`
	When        time.Time `json:"when"`
}

func prepareVm(runtime *goja.Runtime) error {
	if err := runtime.Set("getPrecoBySku", func(sku string) float64 {
		return 1.5
	}); err != nil {
		return err
	}

	if err := runtime.Set("getAliquotaImposto", func(sku string) float64 {
		return 0.275 * float64(len(sku))
	}); err != nil {
		return err
	}

	if err := runtime.Set("getCustoFrete", func(from string, to string, ton float64) float64 {
		return 0.5 * float64(len(from)+len(to)) * ton
	}); err != nil {
		return err
	}

	return nil
}

func run(ctx context.Context) error {

	engHashBytes, err := exec.Command("git", "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		return err
	}

	engHash := strings.TrimSpace(string(engHashBytes))

	h := http.StripPrefix("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		err := func() error {
			vm := goja.New()

			bs, err := os.ReadFile(request.URL.Path)
			if err != nil {
				return err
			}

			if err = prepareVm(vm); err != nil {
				return err
			}

			in := map[string]any{}

			if err = json.NewDecoder(request.Body).Decode(&in); err != nil {
				return err
			}

			if err = vm.Set("$IN", in); err != nil {
				return err
			}

			v, err := vm.RunString(string(bs))
			if err != nil {
				return err
			}

			ret := v.Export()

			fhash := md5.New() // podia ser SHA512, mas sejamos breves nessa demo :)
			fhash.Write(bs)
			ch := fhash.Sum(nil)
			strHash := base64.StdEncoding.EncodeToString(ch)

			out := Out{
				Result:      ret,
				FormulaHash: strHash,
				EngineHash:  engHash,
				Input:       in,
				User:        "getFromIAM",
				When:        time.Now(),
			}

			writer.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(writer).Encode(out)

			return nil
		}()

		if err != nil {
			log.Printf("Error: %s", err.Error())
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

	}))

	return http.ListenAndServe(":8080", h)
}

func main() {
	if err := run(context.Background()); err != nil {
		panic(err)
	}
}
