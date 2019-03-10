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
	Name            string            `json:"name"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func httpClient() *retryablehttp.Client {
	c := retryablehttp.NewClient()
	c.RetryWaitMax = time.Second
	c.RetryMax = 2
	return c
}

func Dir(name, dataDir string) string {
	apiURL := fmt.Sprintf("http://registry.npmjs.org/%s/latest", name)
	resp, err := httpClient().Get(apiURL)
	if err != nil {
		log.Printf("Could not fetch %s: %s", apiURL, err.Error())
		return ""
	}

	var v version
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&v)
	resp.Body.Close()
	if err != nil {
		log.Printf("Could not decode package.json: %s", err.Error())
		return ""
	}

	if v.Repository.Type != "git" || v.Repository.URL == "" {
		return ""
	}

	// download it
	u, err := url.Parse(v.Repository.URL)
	if err != nil {
		log.Printf("Invalid repo URL: %s", err.Error())
		return ""
	}
	p, err := git.Clone(u, dataDir)
	if err != nil {
		log.Printf("Cannot clone repo: %s", err.Error())
		return ""
	}

	return p
}

func Dependencies(repoPath string) []dep.ID {
	// parse the manifest
	pf, err := os.Open(path.Join(repoPath, "package.json"))
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Cannot read package.json: %s", err.Error())
		}
		return nil
	}
	defer pf.Close()

	var pj pkg
	err = json.NewDecoder(pf).Decode(&pj)
	if err != nil {
		log.Printf("Cannot parse package.json: %s", err.Error())
		return nil
	}

	// return dependencies
	var r []dep.ID
	for _, deps := range []map[string]string{pj.Dependencies} { // TODO: add devdependencies
		for d := range deps {
			r = append(r, dep.ID{
				Type: dep.NPM,
				Name: d,
			})
		}
	}
	return r
}
