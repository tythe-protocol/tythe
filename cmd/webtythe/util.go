package main

import (
	"fmt"
	"net/http"
	"os"
)

func badReq(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(msg))
	w.Write([]byte("\n"))
}

func bonk(w http.ResponseWriter, msg string, params ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, params...)
	fmt.Fprintln(os.Stderr)
	w.WriteHeader(http.StatusInternalServerError)
}
