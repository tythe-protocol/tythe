package plan

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// Plan describes the set of packages that should be included in the tythe and
// how the tythe should be distributed amongst them.
type Plan struct {
	Roots   []string           `json:"roots"`
	Weights map[string]float64 `json:"weights"`
}

// AddRoot adds a root package to the plan.
func (c *Plan) AddRoot(r string) {
	for _, ex := range c.Roots {
		if ex == r {
			return
		}
	}
	c.Roots = append(c.Roots, r)
}

// RemoveRoot removes a root package from the plan.
func (c *Plan) RemoveRoot(r string) {
	for i, ex := range c.Roots {
		if ex == r {
			c.Roots = append(c.Roots[0:i], c.Roots[i+1:]...)
			return
		}
	}
}

// SetWeight sets an explicit number of shares for a package in the plan.
func (c *Plan) SetWeight(p string, w float64) {
	if c.Weights == nil {
		c.Weights = map[string]float64{}
	}
	if w < 0.0 {
		panic("s must be positive")
	}
	if w == 1.0 {
		delete(c.Weights, p)
		return
	}
	c.Weights[p] = w
}

// Weight gets the number of shares for a package. If there is no explicit number
// that has been set for the package, the default is 1.0.
func (c *Plan) Weight(p string) float64 {
	if w, ok := c.Weights[p]; ok {
		return w
	}
	return 1.0
}

// Load reads the configuration from a file.
func Load(path string) (Plan, error) {
	w := func(err error) error {
		return errors.Wrapf(err, "Could not load config from: %s", path)
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Plan{}, nil
		}
		return Plan{}, w(err)
	}
	var c Plan
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		return Plan{}, w(err)
	}
	for p, s := range c.Weights {
		if s < 0.0 {
			return Plan{}, w(fmt.Errorf("Negative number of shares not allowed for package: %s: %f", p, s))
		}
	}
	return c, nil
}

// Save writes the configuration to a file.
func Save(c Plan, path string) error {
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
