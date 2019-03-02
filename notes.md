```
# global flags
# The cache dir stores source code and other objects from crawling
# Eventually this should go away since it should be possible to crawl
# git repos somehow directly over http without downloading them
go-tythe [--cache-dir=~/.go-tythe]

go-tythe pay-all <plan>

go-tythe pay-one <project-url> <amount>

# for testing mostly
go-tythe send <address> <money>

# lists all transitive dependencies of <project-url>
go-tythe list
```

todo:

* test `go-tythe list` - created some sample projects, doesn't seem to be working yet
* finish `go-tythe pay-all` - start is in local branch
  - use idempotency token from coinbase api
* implement crawling of npm
* look at shapeshifter to support btc output
* implement coinbase buy widget
