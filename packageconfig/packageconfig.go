package packageconfig

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"

	"github.com/pkg/errors"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// PaymentType is a mechanism that go-tythe can use to move money between parties.
type PaymentType string

const (
	// USDC represents the USDC stablecoin backed by Coinbase and Circle.
	USDC PaymentType = "USDC"
)

var (
	usdcAddressPattern = regexp.MustCompile("^0x[0-9a-f]{40}$")
)

// PackageConfig describes the json metadata developers add to their package to
// opt-in to receiving tythes.
type PackageConfig struct {
	Destination struct {
		Type    PaymentType `json:"type"`
		Address string      `json:"address"`
	} `json:"destination"`
}

func Read(url *url.URL) (PackageConfig, error) {
	var c PackageConfig

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:   url.String(),
		Depth: 1,
	})
	if err == nil {
		h, err := r.ResolveRevision("HEAD")
		if err == nil {
			commit, err := r.CommitObject(*h)
			if err == nil {
				t, err := commit.Tree()
				if err == nil {
					f, err := t.File("tythe.json")
					if err == nil {
						r, err := f.Reader()
						if err == nil {
							d := json.NewDecoder(r)
							err = d.Decode(&c)
						}
					}
				}
			}
		}
	}

	if err != nil {
		return PackageConfig{}, errors.Wrapf(err, "Could not get config for repository: %s", url)
	}

	if c.Destination.Type != USDC {
		return PackageConfig{}, fmt.Errorf("invalid tythe.json - destination type: \"%s\" not supported", c.Destination.Type)
	}

	if !usdcAddressPattern.MatchString(c.Destination.Address) {
		return PackageConfig{}, fmt.Errorf("invalid destination address in tythe.json: \"%s\"", c.Destination.Address)
	}

	return c, nil
}
