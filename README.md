# Overview

The Tythe Protocol is a simple, open, transparent, and license-compatible way to financially support the maintenance of Open Source Software.

# How it Works

### 1. Maintainers add [`tythe.json`](./tythe-sample.json) to the root of their projects

This declares that the maintainer wants to participate, and how to pay them, in a machine-readable format.

### 2. Companies make a public commitment to pay tythes

For example, by posting a link to this page to social media.

### 3. Maintainers get paid

Companies use any compatible [tool](#tools) to [calculate the tythe](#the-tythe-calculation) and distribute it. Companies can divide the tythe however they like, but tools will typically provide sensible defaults.

# The Tythe Calculation

The amount a company should tythe is based on that company's annualized R&D budget. The tythe approaches 1% of this value as the portion of participating dependencies increases.

Specifically:

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


# Tools

* [`go-tythe`](https://github.com/aboodman/go-tythe): The simplest possible CLI client.

# More

* [About](about.md)
* [FAQ](faq.md)
