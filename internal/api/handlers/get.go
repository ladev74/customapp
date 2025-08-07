package handlers

import (
	"fmt"
	"net/http"
)

func Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Get called")
	}
}

func WriteResult() {
}

type response struct {
	Result string `json:"result"`
}
