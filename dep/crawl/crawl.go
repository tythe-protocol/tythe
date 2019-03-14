// Package crawl implements the core of the dependency crawler.
package crawl

import (
	"log"
	"sync"

	"github.com/tythe-protocol/tythe/conf"
	"github.com/tythe-protocol/tythe/dep"
	"github.com/tythe-protocol/tythe/dep/golang"
	"github.com/tythe-protocol/tythe/dep/npm"
	. "github.com/tythe-protocol/tythe/utl/sentinel"
)

// Result contains information about a call to Crawl.
type Result struct {
	Dep      *dep.Dep
	Progress *Progress
}

// Progress indicates the status of the crawl.
type Progress struct {
	Found     int `json:"found"`
	Processed int `json:"processed"`
}

// Crawl performs a parallelized breadth first exploration of the graph rooted at repourl.
// Nodes in the graph are dep.ID, and edges are dependencies to other dep.ID's.
// Child dependencies at each node can be found a variety of strategies, and new strategies
// will be added over time. There will typically be at least one strategy attempted at each
// node per package ecosystem. In the case of Golang, there would eventually be several:
// Go 1.11 modules, Godeps, etc.
func Crawl(repourl, dataDir string, l *log.Logger) <-chan Result {
	const concurrency = 64
	mu := sync.Mutex{}       // protects q, seen, progress
	q := []dep.ID{}          // queue of deps waiting to be explored
	seen := map[dep.ID]SNT{} // deps we've already seen
	progress := Progress{}

	r := make(chan Result)

	// pushes new IDs onto the queue to be processed
	push := func(ids []dep.ID) {
		mu.Lock()
		defer mu.Unlock()
		for _, id := range ids {
			progress.Found++
			q = append(q, id)
		}
		cp := progress
		r <- Result{
			Progress: &cp,
		}
	}

	// pops the next dep or empty string if none left
	pop := func() dep.ID {
		mu.Lock()
		defer mu.Unlock()
		if len(q) == 0 {
			return dep.ID{}
		}
		r := q[0]
		q = q[1:]
		return r
	}

	// marks a dep as having been seen
	// returns true if newly seen, false otherwise
	mark := func(depID dep.ID) bool {
		mu.Lock()
		defer mu.Unlock()
		if _, ok := seen[depID]; ok {
			return false
		}
		seen[depID] = S()
		return true
	}

	// processes a single depname:
	// - fetch the repo
	// - construct the dep for the repo
	// - queue any children for exploration
	processDep := func(depID dep.ID) {
		defer func() {
			mu.Lock()
			progress.Processed++
			cp := progress
			r <- Result{
				Progress: &cp,
			}
			mu.Unlock()
		}()
		if !mark(depID) {
			return
		}
		d, cdns := processDepID(depID, dataDir, l)
		if !d.IsEmpty() {
			r <- Result{
				Dep: &d,
			}
		}
		push(cdns)
	}

	go func() {
		// queue the initial children to explore
		// we don't create a dep for the starting point
		_, cdids := processRepo(repourl, true, l)
		push(cdids)

		// explore the graph concurrently until there are no more depnames queued
		wg := sync.WaitGroup{}
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					n := pop()
					if n.IsEmpty() {
						break
					}
					processDep(n)
				}
			}()
		}
		wg.Wait()
		close(r)
	}()

	return r
}

func processDepID(id dep.ID, dataDir string, l *log.Logger) (r dep.Dep, childDepIDs []dep.ID) {
	l.Printf("Crawling %s", id)

	var dir string
	switch id.Type {
	case dep.NPM:
		dir = npm.Dir(id.Name, dataDir, l)
		break
	case dep.Go:
		dir = golang.Dir(id.Name, dataDir)
		break
	default:
		panic("Unexpected dep type")
	}

	if dir == "" {
		return dep.Dep{}, nil
	}

	c, cdns := processRepo(dir, false, l)
	return dep.Dep{
		ID:   id,
		Conf: c,
	}, cdns
}

func processRepo(path string, isRoot bool, l *log.Logger) (*conf.Config, []dep.ID) {
	opts := npm.Options{
		Dependencies: true,
	}
	if isRoot {
		opts.DevDependencies = true
		opts.PeerDependencies = true
	}
	var r []dep.ID
	ds := npm.Dependencies(path, opts, l)
	r = append(r, ds...)
	ds, err := golang.Dependencies(path)
	if err != nil {
		l.Printf("Could not get golang dependencies: %s", err.Error())
	} else {
		r = append(r, ds...)
	}

	// load conf if any
	c, err := conf.Read(path)
	if err != nil {
		l.Printf("Cannot parse .donate: %s", err.Error())
		return nil, nil
	}

	return c, r
}
