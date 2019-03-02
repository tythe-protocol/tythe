package golang

import (
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasics(t *testing.T) {
	assert := assert.New(t)

	_, thisFile, _, _ := runtime.Caller(0)
	root := path.Dir(path.Dir(path.Dir(thisFile)))
	repoParent := path.Dir(root)
	zTest1 := path.Join(repoParent, "z_test1")
	zTest2 := path.Join(repoParent, "z_test2")

	dirExist := func(p string) bool {
		fi, err := os.Stat(p)
		return err == nil && fi.IsDir()
	}

	if !dirExist(zTest1) || !dirExist(zTest2) {
		assert.Fail("TestList requires supplementary repos zTest1 and zTest2")
	}

	tc := []struct {
		in          string
		expectError bool
		numExpected int
		expectTest2 bool
	}{
		{zTest1, false, 4, true},
		{zTest2, false, 0, false},
		{"not-exist", true, 0, false},
		{root, false, 37, false},
	}

	for _, t := range tc {
		ds, err := List(t.in)

		if t.expectError {
			assert.Error(err)
			assert.Nil(ds)
			continue
		}

		assert.Equal(t.numExpected, len(ds))

		foundTest2 := false
		for _, d := range ds {
			if strings.HasPrefix(d.Name, "github.com/tythe-protocol/z_test2") {
				assert.False(foundTest2)
				assert.Equal("0x1111111111111111111111111111111111111111", d.Conf.USDC)
				foundTest2 = true
			}
		}
		assert.Equal(t.expectTest2, foundTest2)
	}
}
