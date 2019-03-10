package main

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/tythe-protocol/tythe/dep/crawl"
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

		ds := crawl.Crawl(p, cacheDir)

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte{'[', '\n'})

		enc := json.NewEncoder(w)
		first := true
		for d := range ds {
			if first {
				first = false
			} else {
				w.Write([]byte{','})
			}

			enc.Encode(d)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}

		w.Write([]byte{']', '\n'})
	}
}
