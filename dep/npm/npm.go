package npm

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	ppath "path"
	"runtime"
	"strings"
	"sync"

	"github.com/tythe-protocol/go-tythe/conf"
	"github.com/tythe-protocol/go-tythe/dep/shared"

	"github.com/pkg/errors"
)

// List returns all the transitive NPM dependencies of the package at <path>
func List(path string) ([]shared.Dep, error) {
	err := os.Chdir(path)
	if err != nil {
		return nil, err
	}

	defer os.Chdir("-")

	type outchs struct {
		Name string
		Conf *conf.Config
		Err  error
	}
	inch := make(chan string)
	outch := make(chan outchs)

	wg := sync.WaitGroup{}
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			defer wg.Done()
			for in := range inch {
				c, err := readConf(in)
				outch <- outchs{
					Name: in,
					Conf: c,
					Err:  err,
				}
			}
		}()
		wg.Add(1)
	}

	go func() {
		wg.Wait()
		close(outch)
	}()

	go func() {
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

		for p := range pkgs {
			inch <- p
		}

		close(inch)
	}()

	r := []shared.Dep{}
	for out := range outch {
		if out.Err != nil {
			return nil, out.Err
		}
		r = append(r, shared.Dep{
			Name: out.Name,
			Conf: out.Conf,
		})
		fmt.Println(out.Name, out.Conf)
	}

	return r, nil
}

func readConf(name string) (*conf.Config, error) {
	// Seems like there'd be a way to do this from the CLI??
	cmd := exec.Command("node", "-e", fmt.Sprintf("console.log(require.resolve('%s/package.json'))", name))
	out, err := cmd.Output()
	// Have to swallow this error because some packages (e.g., electron) lack a package.json??
	fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error from command: %s", cmd.Args))
	root := ppath.Dir(string(out))
	c, err := conf.Read(root)
	if err != nil {
		return nil, err
	}
	return c, nil
}
