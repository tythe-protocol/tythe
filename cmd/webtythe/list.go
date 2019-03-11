package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
	"github.com/tythe-protocol/tythe/dep/crawl"
	"github.com/tythe-protocol/tythe/git"
)

func list(cacheDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rp := r.URL.Query().Get("r")
		if rp == "" {
			badReq(w, "param r is required")
			return
		}

		repo, err := url.Parse(rp)
		if err != nil {
			badReq(w, errors.Wrapf(err, "invalid parameter r: %s", err).Error())
			return
		}

		if repo.Scheme != "http" && repo.Scheme != "https" {
			badReq(w, "invalid parameter r: scheme must be http or https")
			return
		}
		l := log.New(os.Stdout, "", log.LstdFlags)
		p, err := git.Resolve(repo, cacheDir, l)
		if err != nil {
			bonk(w, err.Error())
			return
		}

		rs := crawl.Crawl(p, cacheDir, l)

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[\n"))

		enc := json.NewEncoder(w)
		for r := range rs {
			if r.Dep != nil {
				enc.Encode(*r.Dep)
			}
			if r.Progress != nil {
				enc.Encode(*r.Progress)
			}
			w.Write([]byte{','})
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}

		w.Write([]byte("null]\n"))
	}
}
