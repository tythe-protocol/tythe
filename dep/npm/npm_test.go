package npm

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchDirs(t *testing.T) {
	assert := assert.New(t)

	ngf := func() []string {
		r, err := nodeGlobalFolders()
		assert.NoError(err)
		return r
	}

	tc := []struct {
		in  string
		exp []string
	}{
		{in: "/foo/bar/baz", exp: append(ngf(), "/foo/bar/baz", "/foo/bar", "/foo")},
		{in: "/foo/node_modules/bar", exp: append(ngf(), "/foo/node_modules/bar", "/foo")},
		{in: "/foo/./bar", exp: append(ngf(), "/foo/bar", "/foo")},
		{in: "/foo/../bar", exp: append(ngf(), "/bar")},
	}

	for _, t := range tc {
		act, err := searchDirs(t.in)
		assert.NoError(err)
		assert.Equal(t.exp, act)
	}

	wd, err := os.Getwd()
	assert.NoError(err)

	tc2 := []struct {
		in        string
		expSameAs string
	}{
		{in: ".", expSameAs: wd},
		{in: "./foo", expSameAs: path.Join(wd, "foo")},
		{in: "foo", expSameAs: path.Join(wd, "foo")},
	}

	for _, t := range tc2 {
		exp, err := searchDirs(t.expSameAs)
		assert.NoError(err)
		act, err := searchDirs(t.in)
		assert.NoError(err)
		assert.Equal(exp, act)
	}
}
