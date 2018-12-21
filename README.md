# Tythe

1% of R&D to Open Source Maintenance

# What

Tythe automatically directs 1% of your company's R&D budget to the maintenance of the open source projects you depend on.

# How

1. Open source maintainers add [tythe.json](./tythe-sample.json) to their repositories. This declares how to send them money in a machine-readable way.
2. Companies take [The Tythe Covenant](./covenant.md) by posting it to social media, or on their website. The Covenant is a public promise to contribute [up to 1%](#how-tythes-are-calculated) of R&D monthly to open source maintenance. Enforcement of the convenant is entirely social.
3. Companies use [go-tythe](#status) (or whatever other tool they want) to automatically distribute and send money to the maintainers of their dependencies every month.
4. 🙌

# Why

Open Source started out as a fringe movement, but over several decades, it has grown into critical shared infrastructure. To a close approximation, *all* companies are now heavily dependent upon open source. We won!

However. Open source is also largely unmaintained. The people most capable of doing the maintenance are usually doing so at night and on weekends, while they work on something unrelated during the day to pay the bills.

It’s time to evolve. We need to direct resources to the maintenance of our digital commons, and allow the right people to do that full-time. Tythe is one easy way to do this.

# The Tythe Calculation

A value of a company's tythe is based on its annualized R&D expenditure. It will never be greater than 1% of this value, and usually significantly smaller.

```
tythe = R&D * 0.01 * tythed_deps / total_deps

R&D:         annualized R&D expenditure
total_deps:  count of transitive dependencies
tythed_deps: count of total_deps that contain a tythe.json
```

## Example 1

 * Your current R&D expenditure: `$2M/yr`
 * Number of transitive dependencies in your tree: `500`
 * Number of transitive dependencies that include `tythe.json`: `150`
 
Your tythe is: `$2M * 0.01 * 150 / 500 = $6000/year` or `$500/month` or `$3.33/mo/dep`

## Example 2

 * Your current R&D expenditure: `$16B/yr`
 * Number of transitive dependencies in your tree: `10k`
 * Number of transitive dependencies that include `tythe.json`: `2k`

Your tythe is: `$16B * 0.01 * 2000 / 10000 = $32M/year` or `$2.7M/month` or about `$1.3k/mo/dep`

# Make the Tythe Covenant

Do you run a company that uses open source? You should make the Tythe Covenant, documenting your commitment to support open source maintenance.

# Status

Right now this document is all there is. It’s in the collecting-feedback stage.

Given interest, I eventually imagine a series of tools that plug into continuous integration that companies can use to calculate and pay tythes automatically.

# FAQ

### What's with the name?
A tithe is a donation that many religious organizations ask of their practicioners, paid to support the clergy and other shared infrastructure of the church. It comes from the old english "tenth", because traditionally the tythe is 10% of ones income.

### Why “tythe” not “tithe”?
It’s the archaic spelling. I just like it better.

### Why is the tythe based on R&D, not revenue or profit?
 * Revenue scales too fast in most companies, it ends up being impractical at either the high or low end.
 * R&D scales sub-linear with revenue, making it a nice fit for this use case.
 * R&D in software companies is almost entirely engineers. This makes it easy to compare to the value that open source provides.
 * Lots of companies have no revenue, but are funded. They should still contribute.
 * R&D is already reported publicly in many companies, leading to transparency.
