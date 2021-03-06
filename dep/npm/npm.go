// Package npm implements the npm-specific bits of dependency crawling.
package npm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/tythe-protocol/tythe/dep"
	"github.com/tythe-protocol/tythe/git"
)

type version struct {
	Name       string     `json:"name"`
	Repository repository `json:"repository"`
}

type repository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type dist struct {
	Shasum  string `json:"shasum"`
	Tarball string `json:"tarball"`
}

type pkg struct {
	Name             string            `json:"name"`
	Dependencies     map[string]string `json:"dependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
	PeerDependencies map[string]string `json:"peerDependencies"`
}

func httpClient(l *log.Logger) *retryablehttp.Client {
	c := retryablehttp.NewClient()
	c.RetryWaitMax = time.Second
	c.RetryMax = 2
	c.Logger = l
	return c
}

func Dir(name, dataDir string, l *log.Logger) string {
	apiURL := fmt.Sprintf("http://registry.npmjs.org/%s/latest", name)
	resp, err := httpClient(l).Get(apiURL)
	if err != nil {
		l.Printf("Could not fetch %s: %s", apiURL, err.Error())
		return ""
	}

	var v version
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&v)
	resp.Body.Close()
	if err != nil {
		l.Printf("Could not decode package.json: %s", err.Error())
		return ""
	}

	if v.Repository.Type != "git" || v.Repository.URL == "" {
		return ""
	}

	// download it
	u, err := url.Parse(v.Repository.URL)
	if err != nil {
		l.Printf("Invalid repo URL: %s", err.Error())
		return ""
	}
	p, err := git.Clone(u, dataDir, l)
	if err != nil {
		l.Printf("Cannot clone repo: %s", err.Error())
		return ""
	}

	return p
}

type Options struct {
	Dependencies     bool
	DevDependencies  bool
	PeerDependencies bool
}

func Dependencies(repoPath string, opts Options, l *log.Logger) []dep.ID {
	// parse the manifest
	pf, err := os.Open(path.Join(repoPath, "package.json"))
	if err != nil {
		if !os.IsNotExist(err) {
			l.Printf("Cannot read package.json: %s", err.Error())
		}
		return nil
	}
	defer pf.Close()

	var pj pkg
	err = json.NewDecoder(pf).Decode(&pj)
	if err != nil {
		l.Printf("Cannot parse package.json: %s", err.Error())
		return nil
	}

	// return dependencies
	depss := []map[string]string{}
	if opts.Dependencies {
		depss = append(depss, pj.Dependencies)
	}
	if opts.DevDependencies {
		depss = append(depss, pj.DevDependencies)
	}
	if opts.PeerDependencies {
		depss = append(depss, pj.PeerDependencies)
	}
	var r []dep.ID
	for _, deps := range depss {
		for d := range deps {
			r = append(r, dep.ID{
				Type: dep.NPM,
				Name: d,
			})
		}
	}
	return r
}
