package motor_calculo

import (
	"encoding/json"
	"net/http"
)

type Api struct {
	Service *Service
}

func (a *Api) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	err := func() error {
		defer request.Body.Close()
		in := &Request{}

		if err := json.NewDecoder(request.Body).Decode(in); err != nil {
			return err
		}

		res, err := a.Service.Calc(request.Context(), in)
		if err != nil {
			return err
		}

		return json.NewEncoder(writer).Encode(res)
	}()

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
