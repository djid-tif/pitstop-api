package utils

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func ExtractUintFromRequest(name string, r *http.Request) (uint, error) {
	vars := mux.Vars(r)

	val, found := vars[name]
	if !found {
		return 0, fmt.Errorf("%s not provided", name)
	}

	res, err := StringToUint(val)
	if err != nil {
		return 0, fmt.Errorf("%s is invalid", name)
	}

	return res, nil
}
