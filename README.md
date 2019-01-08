# Overview

This document describes a simple, open, transparent, and license-comaptible way to financially support improved maintenance of Open Source software. Discussion and suggestions for improvement are requested.

# Problem

Over the past several decades, Open Source has gone from a fringe hippie movement to critical global infrastructure.

We won! Basically all companies are now heavily dependent on Open Source software. And the spoils are rich — a vast and deep software commons which enable us to build at a previously unimaginable scale and pace.

Unfortunately, this common infrastructure is also chronically under-funded. You can't make a living maintaining popular and widely used open source packages, so developers typically do it in their spare time.

One result has been a steady drip of major security issues. But more importantly, the best work just isn't getting done. The people suited to tending our digital commons are building better ad targeting to pay the bills.

# Proposal

An *Open Source Tythe* is up to 1% of a company's R&D budget, distributed continuously amongst the maintainers of the Open Source projects the company depends on.

All companies that use Open Source are encouraged to participate in the Tythe Protocol.

It works like this:

1. Open source maintainers add [`tythe.json`](./tythe-sample.json) to their repositories. This declares that the developer wants to participate in tythe, and how they should get paid, in a machine-readable way. It also contains a copy of The Tythe Covenant - which describes a minimum level of craftsmanship and responsibility for the library that the developer commits to.
2. Companies make a public statement committing to the Tythe Protocol (e.g., by posting on social media).
3. Every month companies calculate their tythe and distribute it amongst their open source dependencies in any way they see fit. All dependencies, transitively, are elibile, including those that are used via source code, object code, or remote interface. Dependencies can be discovered via `tythe.json` using any tool, or alternately by manual specification.
4. 🙌

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


# Dividing the Tythe

How a company divides its tythe amongst its dependencies is entirely up to that company. The only requirement from the protocol is the total amount of the tythe, not to where it goes.

Tools will provide a variety of features and options for dividing the tythe.

# Tools

Currently, nothing exists, but work will begin soon on `go-tythe`, the simplest possible command-line implementation of a client that can distribute payments once.

# FAQ

### Why “tythe” not “tithe”?
It’s the archaic spelling. I just like it better.

### Why is the tythe based on R&D, not revenue or profit?
 * Revenue scales too fast in most companies, it ends up being impractical at either the high or low end.
 * R&D scales sub-linear with revenue, making it a nice fit for this use case.
 * Lots of companies have no revenue, but are funded. They should still contribute.
 * R&D in software companies is almost entirely engineers. This makes it easy to compare to the value that open source provides.
 * R&D is already reported publicly in many companies, leading to transparency.
