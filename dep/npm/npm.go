package npm

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	ppath "path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tythe-protocol/tythe/conf"
	"github.com/tythe-protocol/tythe/dep/shared"

	homedir "github.com/mitchellh/go-homedir"
)

// List returns all the transitive NPM dependencies of the package at <path>
func List(path string) ([]shared.Dep, error) {
	err := os.Chdir(path)
	if err != nil {
		return nil, err
	}

	defer os.Chdir("-")

	_, err = os.Stat("package.json")
	if os.IsNotExist(err) {
		return nil, nil
	}

	// Ignore errors because ls returns errors for unmet peer deps, even when it is still returning useful info.
	out, _ := exec.Command("npm", "ls").Output()
	s := bufio.NewScanner(bytes.NewReader(out))
	pkgs := map[string]struct{}{}

	// Skip the first line - it is the root module itself
	s.Scan()
	for s.Scan() {
		t := s.Text()
		for _, f := range strings.Split(t, " ") {
			if strings.Contains(f, "@") {
				p := f[:strings.LastIndex(f, "@")]
				pkgs[p] = struct{}{}
			}
		}
	}

	// This is basically an implementation of require.resolve, from:
	// https://nodejs.org/api/modules.html#modules_all_together
	dirs, err := searchDirs(path)
	if err != nil {
		return nil, err
	}

	r := []shared.Dep{}
	for pkg := range pkgs {
		for _, dir := range dirs {
			modpath := ppath.Join(dir, "node_modules", pkg)
			_, err := os.Stat(modpath)
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return nil, err
			}
			c, err := conf.Read(modpath)
			if err != nil {
				return nil, err
			}
			d := shared.Dep{
				Name: pkg,
				Conf: c,
			}
			r = append(r, d)
			break
		}
	}
	return r, nil
}

func searchDirs(path string) ([]string, error) {
	dirs, err := nodeGlobalFolders()
	if err != nil {
		return nil, err
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	for path != "" {
		dir, leaf := ppath.Split(path)
		if len(dir) > 0 {
			dir = dir[:len(dir)-1]
		}
		if leaf != "node_modules" {
			dirs = append(dirs, path)
		}
		path = dir
	}

	return dirs, nil
}

func nodeGlobalFolders() ([]string, error) {
	dirs := nodePath()
	hd, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	dirs = append(dirs, ppath.Join(hd, ".node_modules"))
	dirs = append(dirs, ppath.Join(hd, ".node_libraries"))
	// TODO: node_prefix - I don't get this bit
	return dirs, nil
}

func nodePath() []string {
	np := os.Getenv("NODE_PATH")
	if np == "" {
		return nil
	}
	sep := ":"
	if runtime.GOOS == "windows" {
		sep = ";"
	}
	return strings.Split(np, sep)
}
