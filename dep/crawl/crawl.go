package crawl

import (
	"log"
	"sync"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"

	"github.com/tythe-protocol/tythe/conf"
	"github.com/tythe-protocol/tythe/dep"
	"github.com/tythe-protocol/tythe/dep/golang"
	"github.com/tythe-protocol/tythe/dep/npm"
	. "github.com/tythe-protocol/tythe/utl/sentinel"
)

func httpClient() *retryablehttp.Client {
	c := retryablehttp.NewClient()
	c.RetryWaitMax = time.Second
	c.RetryMax = 2
	return c
}

func Crawl(repourl, dataDir string) []dep.Dep {
	// Crawl performs a parallelized breadth first exploration of the graph rooted at repourl
	// Nodes in the graph are dep.ID, and edges are child dependencies represented as (Dep)
	// Child dependencies at each node can be found a variety of ways, and this will improve over time
	// There will typically be at least one strategy attempted at each node per package ecosystem
	// In the case of Golang, there would eventually be several: Go 1.11 modules, Godeps, etc.

	const concurrency = 64
	q := []dep.ID{}          // queue of deps waiting to be explored
	seen := map[dep.ID]SNT{} // deps we've already seen
	r := []dep.Dep{}         // deps we have found
	mu := sync.Mutex{}       // protects q, seen, r

	// pushes new IDs onto the queue to be processed
	push := func(depIDs []dep.ID) {
		mu.Lock()
		defer mu.Unlock()
		q = append(q, depIDs...)
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

	// adds a new dep to the resultset
	collect := func(d dep.Dep) {
		mu.Lock()
		defer mu.Unlock()
		r = append(r, d)
	}

	// processes a single depname:
	// - fetch the repo
	// - construct the dep for the repo
	// - queue any children for exploration
	processDep := func(depID dep.ID) {
		if !mark(depID) {
			return
		}
		d, cdns := processDepID(depID, dataDir)
		if !d.IsEmpty() {
			collect(d)
		}
		push(cdns)
	}

	// queue the initial children to explore
	// we don't create a dep for the starting point
	_, cdids := processRepo(repourl)
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

	return r
}

func processDepID(id dep.ID, dataDir string) (r dep.Dep, childDepIDs []dep.ID) {
	log.Printf("Crawling %s", id)

	var dir string
	switch id.Type {
	case dep.NPM:
		dir = npm.Dir(id.Name, dataDir)
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

	c, cdns := processRepo(dir)
	return dep.Dep{
		ID:   id,
		Conf: c,
	}, cdns
}

func processRepo(path string) (*conf.Config, []dep.ID) {
	var r []dep.ID
	ds := npm.Dependencies(path)
	r = append(r, ds...)
	ds, err := golang.Dependencies(path)
	if err != nil {
		log.Printf("Could not get golang dependencies: %s", err.Error())
	} else {
		r = append(r, ds...)
	}

	// load conf if any
	c, err := conf.Read(path)
	if err != nil {
		log.Printf("Cannot parse .donate: %s", err.Error())
		return nil, nil
	}

	return c, r
}
