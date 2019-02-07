
```
# global flags
# You could have different configs for different top-level entities
# Seems unusual though
go-tythe [--config=~/.go-tythe]

# dep subcommand manages roots in the dependency graph
# can be either your own repos, or else other top-level deps
# that aren't discoverable, e.g., linux
go-tythe dep add <package>
go-tythe dep remove

# shows transitive dep tree
go-tythe dep list

# show or set an explicit weight for a package (default 1)
go-tythe weight <package> [<weight>]

# serves a UI for managing weights?? (or maybe UI is a separate layer)
# go-tythe ui

# crawls dependency tree starting from roots
# calculates tythes from R&D <basis>
# presents final invoice and asks for confirmation
# this would store state in config so that it can be idempotent
# if run twice in a row, second one would be no-op until one month has elapsed
go-tythe run <basis>
```

Problem:
- Don't want it running unmonitored ... needs to prompt user somehow
- In a server world, I imagine this sending an email to an admin that links to a UI
- In a desktop world, it could be a notification