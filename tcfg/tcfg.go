package tcfg

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// Config is the config data for the go-tythe program.
type Config struct {
	Roots   []string           `json:"roots"`
	Weights map[string]float64 `json:"weights"`
}

// AddRoot adds a root to the dependency graph.
func (c *Config) AddRoot(r string) {
	for _, ex := range c.Roots {
		if ex == r {
			return
		}
	}
	c.Roots = append(c.Roots, r)
}

// RemoveRoot removes a root from the dependency graph.
func (c *Config) RemoveRoot(r string) {
	for i, ex := range c.Roots {
		if ex == r {
			c.Roots = append(c.Roots[0:i], c.Roots[i+1:]...)
			return
		}
	}
}

// SetWeight sets an explicit number of shares for a package in the dependency graph.
func (c *Config) SetWeight(p string, s float64) {
	if c.Weights == nil {
		c.Weights = map[string]float64{}
	}
	if s < 0.0 {
		panic("s must be positive")
	}
	c.Weights[p] = s
}

// Weight gets the number of shares for a package. If there is no explicit number
// that has been set for the package, the default is 1.0.
func (c *Config) Weight(p string) float64 {
	if w, ok := c.Weights[p]; ok {
		return w
	}
	return 1.0
}

// Load reads the configuration from a file.
func Load(path string) (Config, error) {
	w := func(err error) error {
		return errors.Wrapf(err, "Could not load config from: %s", path)
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, w(err)
	}
	var c Config
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return Config{}, w(err)
	}
	for p, s := range c.Weights {
		if s < 0.0 {
			return Config{}, w(fmt.Errorf("Negative number of shares not allowed for package: %s: %f", p, s))
		}
	}
	return c, nil
}

// Save writes the configuration to a file.
func Save(c Config, path string) error {
	w := func(err error) error {
		return errors.Wrapf(err, "Could not save config: %+v to: %s", c, path)
	}
	f, err := os.Create(path)
	if err != nil {
		return w(err)
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(c)
	if err != nil {
		return w(err)
	}

	return nil
}
