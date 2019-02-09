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
