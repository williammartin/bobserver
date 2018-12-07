package main

import (
	"fmt"
	"net/http"

	"github.com/piquette/finance-go/equity"
)

func main() {
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/fetch", fetch)
	if err := http.ListenAndServe(":9000", nil); err != nil {
		panic(err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func fetch(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	ticker := queryParams.Get("ticker")

	if ticker == "" {
		w.WriteHeader(400)
		w.Write([]byte("please provide a ticker query parameter"))
		return
	}

	quote, err := equity.Get(ticker)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("an error occurred fetching the price"))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("%f", quote.RegularMarketPrice)))
}
