# Tythe

1% of R&D for the open source projects you depend on.

# What

Tythe automatically directs 1% of your company's R&D budget to the maintenance of the open source projects you depend on:

1. Maintainers add the [tythe.json](./tythe-sample.json) file to their repositories. This declares how to send them money in a machine-readable way.
2. Companies install and run Tythe on their own servers.
3. Tythe monitors the company's dependency tree and automatically pays the maintainers every month. By default the tythe is split evenly amongst all dependencies, but users can adjust the share if desired.

# Why

Open Source started out as a fringe movement, but over several decades, it has grown into absolutely critical shared infrastructure. To a close approximation, *all* companies are now heavily dependent upon open source. We won!

However. Open source is also basically unmaintained. The people most capable of doing the maintenance are typically doing so at night and on weekends, while they work on something unrelated during the day to make a living.

It’s time to evolve. We need to direct resources to maintaining our digital commons and allow the right people to do that full-time.

# How the Tythe is Calculated

The expected tythe is based on a company's annualized R&D expenditure. It will never be greater than 1% of this value, and usually significantly smaller.

```
tythe = R&D * 0.1 * tythed_deps / total_deps

R&D:         annualized R&D expenditure
total_deps:  count of transitive dependencies
tythed_deps: count of total_deps that contain a tythe.json
```

## Example 1

 * Your current R&D expenditure: `$2M/yr`
 * Number of transitive dependencies in your tree: `500`
 * Number of transitive dependencies that include `tythe.json`: `150`
 
Your tythe is: `$2M * 0.1 * 150 / 500 = $6000/year` or `$500/month`

## Example 2

 * Your current R&D expenditure: `$16B/yr`
 * Number of transitive dependencies in your tree: `10k`
 * Number of transitive dependencies that include `tythe.json`: `2k`

Your tythe is: `$16B * 0.1 * 2000 / 10000 = $32M/year` or `$2.7M/month` or about `$1.3k/mo/dep`

# Make the Tythe Covenant

Do you run a company that uses open source? You should make the Tythe Covenant, documenting your commitment to support open source maintenance.

# Status

Right now this document is all there is. It’s in the collecting-feedback stage.

Given interest, I eventually imagine a series of tools that plug into continuous integration that companies can use to calculate and pay tythes automatically.

# FAQ

### Why “tythe” not “tithe”?
It’s the archaic spelling. I just like it better.

### Why is the tythe based on R&D, not revenue or profit?
 * Revenue scales too fast in most companies, it ends up being impractical at either the high or low end.
 * R&D scales sub-linear with revenue, making it a nice fit for this use case.
 * Lots of companies have no revenue, but are funded. They should still contribute.
 * R&D in software companies is almost entirely engineers. This makes it easy to compare to the value that open source provides.
 * R&D is already reported publicly in many companies, leading to transparency.
