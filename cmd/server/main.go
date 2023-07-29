package main

import (
	"net/http"
)

const addr = `:8080`

func updateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	http.HandleFunc(`/update`, updateMetric)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
