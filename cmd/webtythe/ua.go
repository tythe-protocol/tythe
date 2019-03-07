package main

import (
	"fmt"
	"net/http"

	"github.com/attic-labs/noms/go/d"

	"github.com/ua-parser/uap-go/uaparser"
)

func ua(w http.ResponseWriter, r *http.Request) {
	parser, err := uaparser.NewFromBytes(regexes)
	d.PanicIfError(err)

	client := parser.Parse(r.UserAgent())
	w.Write([]byte(fmt.Sprintf("%+v", *client.UserAgent)))
}
