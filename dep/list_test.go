package dep

import (
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"runtime"
	"testing"

	// TODO: Don't really like this check package - either get or recreate the one from PL
	chk "gopkg.in/check.v1"

	"github.com/tythe-protocol/go-tythe/conf"
	"github.com/tythe-protocol/go-tythe/git"
)

func Test(t *testing.T) { chk.TestingT(t) }

type DepSuite struct{}

var _ = chk.Suite(&DepSuite{})

func (s *DepSuite) TestList(c *chk.C) {
	_, thisFile, _, _ := runtime.Caller(0)

	root := path.Dir(path.Dir(thisFile))
	repoParent := path.Dir(root)
	zTest1 := path.Join(repoParent, "z_test1")
	zTest2 := path.Join(repoParent, "z_test2")

	dirExist := func(p string) bool {
		fi, err := os.Stat(p)
		return err == nil && fi.IsDir()
	}

	if !dirExist(zTest1) || !dirExist(zTest2) {
		c.Error("TestList requires supplementary repos zTest1 and zTest2")
	}

	tc := []struct {
		in          string
		expectError bool
		expectClone bool
		numExpected int
		expectTest2 bool
	}{
		{zTest1, false, false, 1, true},
		{zTest2, false, false, 0, false},
		{"not-exist", true, false, 0, false},
		{"file:" + zTest1, false, true, 1, true},
		{"file:" + zTest2, false, true, 0, false},
		{"file:not-exist", true, false, 0, false},
		{root, false, false, 37, false},
	}

	for _, t := range tc {
		dataDir, err := ioutil.TempDir("", "")
		c.Assert(err, chk.IsNil)
		u, err := url.Parse(t.in)
		c.Assert(u, chk.NotNil)
		c.Assert(err, chk.IsNil)
		ds, err := List(u, dataDir)

		if t.expectError {
			c.Check(err, chk.NotNil)
			c.Check(ds, chk.IsNil)
			continue
		}

		if t.expectClone {
			d, err := os.Open(git.DirForURL(u, dataDir))
			c.Assert(err, chk.IsNil)
			ns, err := d.Readdirnames(-1)
			c.Assert(err, chk.IsNil)
			c.Check(ns, chk.Not(chk.HasLen), 0)
		} else {
			_, err := os.Stat(git.DirForURL(u, dataDir))
			if err == nil || !os.IsNotExist(err) {
				panic("expected dir to be empty")
			}
		}

		c.Assert(ds, chk.HasLen, t.numExpected)

		foundTest2 := false
		for _, d := range ds {
			if d.Name == "github.com/tythe-protocol/z_test2@v0.0.0-20190209085012-7a77ae91ad6e" {
				c.Check(foundTest2, chk.Not(chk.Equals), true)
				c.Check(d.Conf.Destination.Address, chk.Equals, "0x1111111111111111111111111111111111111111")
				c.Check(d.Conf.Destination.Type, chk.Equals, conf.USDC)
				foundTest2 = true
			}
		}
		c.Check(foundTest2, chk.Equals, t.expectTest2)
	}
}
