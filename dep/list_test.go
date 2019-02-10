package dep

import (
	"os"
	"path"
	"runtime"
	"testing"

	// TODO: Don't really like this check package - either get or recreate the one from PL
	chk "gopkg.in/check.v1"

	"github.com/tythe-protocol/go-tythe/conf"
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
		numExpected int
		expectTest2 bool
	}{
		{zTest1, false, 1, true},
		{zTest2, false, 0, false},
		{"not-exist", true, 0, false},
		{root, false, 37, false},
	}

	for _, t := range tc {
		ds, err := List(t.in)

		if t.expectError {
			c.Check(err, chk.NotNil)
			c.Check(ds, chk.IsNil)
			continue
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
