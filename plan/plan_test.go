package plan

import (
	"io/ioutil"
	"path"
	"testing"

	// TODO: Don't really like this check package - either get or recreate the one from PL
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type PlanSuite struct{}

var _ = Suite(&PlanSuite{})

func (s *PlanSuite) TestLoadNoFile(c *C) {
	p, err := ioutil.TempDir("", "")
	c.Assert(err, IsNil)
	p = path.Join(p, "foo")
	cfg, err := Load(p)
	c.Assert(err, IsNil)
	c.Check(cfg, DeepEquals, Plan{})
}

func (s *PlanSuite) TestLoadRoundTrip(c *C) {
	f, err := ioutil.TempFile("", "")
	c.Assert(err, IsNil)
	cfg := Plan{}
	cfg.Roots = []string{"foo", "bar"}
	cfg.Weights = map[string]float64{"baz": 2, "hotdog": 0.5}
	err = Save(cfg, f.Name())
	c.Assert(err, IsNil)

	cfg2, err := Load(f.Name())
	c.Assert(err, IsNil)
	c.Check(cfg, DeepEquals, cfg2)
}

func (s *PlanSuite) TestRoots(c *C) {
	cfg := Plan{}
	cfg.AddRoot("foo")
	cfg.AddRoot("bar")
	cfg.AddRoot("foo")
	c.Check(cfg, DeepEquals, Plan{Roots: []string{"foo", "bar"}})

	cfg.RemoveRoot("bar")
	cfg.RemoveRoot("baz")
	c.Check(cfg, DeepEquals, Plan{Roots: []string{"foo"}})

	cfg.RemoveRoot("foo")
	c.Check(cfg, DeepEquals, Plan{Roots: []string{}})
}

func (s *PlanSuite) TestWeights(c *C) {
	cfg := Plan{}
	cfg.SetWeight("foo", 2)
	cfg.SetWeight("foo", 3)
	cfg.SetWeight("bar", 0.5)
	c.Check(cfg, DeepEquals, Plan{Weights: map[string]float64{"foo": 3, "bar": 0.5}})
	c.Check(cfg.Weight("foo"), Equals, 3.0)
	c.Check(cfg.Weight("bar"), Equals, 0.5)
}
