package main

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/attic-labs/noms/go/d"

	"github.com/pkg/errors"
	"github.com/tythe-protocol/tythe/dep"
	"github.com/tythe-protocol/tythe/git"
)

func list(cacheDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rs := r.URL.Query().Get("r")
		if rs == "" {
			badReq(w, "param r is required")
			return
		}

		repo, err := url.Parse(rs)
		if err != nil {
			badReq(w, errors.Wrapf(err, "invalid parameter r: %s", err).Error())
			return
		}

		if repo.Scheme != "http" && repo.Scheme != "https" {
			badReq(w, "invalid parameter r: scheme must be http or https")
			return
		}

		p, err := git.Resolve(repo, cacheDir)
		if err != nil {
			bonk(w, err.Error())
			return
		}

		ds, err := dep.List(p)
		if err != nil {
			bonk(w, err.Error())
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(ds)
		d.PanicIfError(err)
	}
}
