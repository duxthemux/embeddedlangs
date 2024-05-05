package motor_calculo

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	sjson "go.starlark.net/lib/json"
	"go.starlark.net/lib/math"
	"go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"
)

type Service struct {
	Path string
}

type PrdQntEntry struct {
	SKU  string
	Qntd string
}

type Request struct {
	Src             string             `json:"src"`
	IDCliente       string             `json:"id_cliente"`
	Produtos        map[string]float64 `json:"produtos"`
	EnderecoEntrega string             `json:"endereco_entrega"`
	Estado          string             `json:"estado"`
}

func (r *Request) AsStarlark() (*starlark.Dict, error) {
	ret := starlark.NewDict(0)
	for _, err := range []error{
		ret.SetKey(starlark.String("src"), starlark.String(r.Src)),
		ret.SetKey(starlark.String("id_cliente"), starlark.String(r.IDCliente)),
		ret.SetKey(starlark.String("endereco_entrega"), starlark.String(r.EnderecoEntrega)),
		ret.SetKey(starlark.String("estado"), starlark.String(r.Estado)),
	} {
		if err != nil {
			return nil, err
		}
	}

	produtosDict := starlark.NewDict(0)
	for sku, qntd := range r.Produtos {
		if err := produtosDict.SetKey(starlark.String(sku), starlark.Float(qntd)); err != nil {
			return nil, err
		}
	}

	if err := ret.SetKey(starlark.String("produtos"), produtosDict); err != nil {
		return nil, err
	}

	return ret, nil
}

type Response struct {
	Result any    `json:"result"`
	Hash   string `json:"hash"`
}

func getHash() (string, error) {
	engHashBytes, err := exec.Command("git", "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(engHashBytes)), nil
}

func (s *Service) Init() error {
	var err error

	s.Path, err = filepath.Abs(s.Path)
	if err != nil {
		return err
	}

	if err := os.Chdir(s.Path); err != nil {
		return err
	}

	hash, err := getHash()
	if err != nil {
		return err
	}
	slog.Info("motor_calculo.Service.Init ok", "hash:", hash, "wd", s.Path)
	return nil
}

func PrecoSku(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.Float(2.0), nil
}

func DescontoCliente(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	cliente := args.Index(0).(starlark.String)
	switch cliente {
	case "001":
		return starlark.Float(0.03), nil
	case "002":
		return starlark.Float(0.02), nil
	case "003":
		return starlark.Float(0.01), nil
	default:
		return starlark.Float(0.0), nil
	}

}

func TributacaoEstado(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	estado := args.Index(0).(starlark.String)

	switch estado {
	case "SC":
		return starlark.Float(0.1), nil
	case "MG":
		return starlark.Float(0.15), nil
	case "SP":
		return starlark.Float(0.25), nil
	default:
		return starlark.Float(0.35), nil
	}

	return starlark.Float(2.0), nil
}

var Module = &starlarkstruct.Module{
	Name: "motor",
	Members: starlark.StringDict{
		"preco_sku":         starlark.NewBuiltin("preco_sku", PrecoSku),
		"desconto_cliente":  starlark.NewBuiltin("desconto_cliente", DescontoCliente),
		"tributacao_estado": starlark.NewBuiltin("tributacao_estado", TributacaoEstado),
	},
}

func DictToMap(d *starlark.Dict) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	for _, k := range d.Keys() {
		v, ok, _ := d.Get(k)
		if ok {
			switch v := v.(type) {
			case *starlark.List:
				listVal := []any{}
				for i := 0; i < v.Len(); i++ {
					listVal = append(listVal, v.Index(i).String())
				}
				ret[k.String()] = listVal
			case *starlark.Dict:
				nv, err := DictToMap(v)
				if err != nil {
					return nil, err
				}
				ret[k.String()] = nv
			default:
				ret[k.String()] = v.String()
			}
		}
	}

	return ret, nil
}

func (s *Service) Calc(ctx context.Context, request *Request) (*Response, error) {
	hash, err := getHash()
	if err != nil {
		return nil, err
	}

	starlark.Universe["motor"] = Module
	starlark.Universe["json"] = sjson.Module
	starlark.Universe["time"] = time.Module
	starlark.Universe["math"] = math.Module

	thread := &starlark.Thread{
		Name: "my thread",
		Print: func(_ *starlark.Thread, msg string) {
			fmt.Println(msg)
		},
	}

	srcFullPath := filepath.Join(s.Path, request.Src)

	srcBs, err := os.ReadFile(srcFullPath)
	if err != nil {
		return nil, err
	}

	_, prog, err := starlark.SourceProgramOptions(&syntax.FileOptions{}, "main.star", string(srcBs), func(s string) bool {
		switch s {
		case "ctx":
			return true

		}
		return false
	})

	if err != nil {
		return nil, err
	}

	startlarkRequest, err := request.AsStarlark()
	if err != nil {
		return nil, err
	}
	ret, err := prog.Init(thread, starlark.StringDict{
		"ctx": startlarkRequest,
	})

	dRet := ret["ret"].(*starlark.Dict)

	retGo, err := DictToMap(dRet)
	if err != nil {
		return nil, err
	}

	bs, err := json.Marshal(ret)
	nret := map[string]any{}
	json.Unmarshal(bs, &nret)

	return &Response{Result: retGo, Hash: hash}, err
}
