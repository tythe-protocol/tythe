# Overview

This document describes a simple, open, transparent, and license-compatible way to financially support improved maintenance of Open Source software. Discussion and suggestions for improvement are requested.

# Problem

Over the past several decades, Open Source has gone from a fringe hippie movement to critical global infrastructure.

We won! Basically all companies are now heavily dependent on Open Source software. And the spoils are rich — a vast and deep software commons which enable us to build at a previously unimaginable scale and pace.

Unfortunately, this common infrastructure is also chronically under-funded. You can't make a living maintaining popular and widely used open source packages, so developers typically do it in their spare time.

One result has been a steady drip of major security issues. But more importantly, the best work just isn't getting done. The people suited to tending our digital commons are building better ad targeting (or whatever) to pay the bills.

# Proposal

An *Open Source Tythe* is up to 1% of a company's R&D budget, distributed continuously amongst the maintainers of the Open Source projects the company depends on. All companies that use Open Source are encouraged to participate.

It works like this:

1. Open source maintainers add [`tythe.json`](./tythe-sample.json) to their repositories. This declares that the developer wants to participate, and how they should get paid, in a machine-readable way. It also contains a copy of *The Tythe Covenant* - which describes a minimum level of craftsmanship and responsibility for the library that the developer commits to.
2. Companies make a public statement committing to the Tythe Protocol (e.g., by posting on social media).
3. Every month companies calculate their tythe and distribute it amongst their dependencies in any way they see fit. All dependencies, transitively, are elibile, including those that are used via source code, object code, or remote interface. Dependencies can be discovered via `tythe.json` using any tool, or alternately by manual specification.
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

### What's with the name?

A tithe is a donation that many religious organizations ask of their practicioners, paid to support the clergy and other shared infrastructure of the church. It comes from the old english "tenth", because traditionally it was 10% of one's income.

### Why “tythe” not “tithe”?
It’s the archaic spelling. I just like it better.

### Why 1%?
There should be some amount that is easy for companies to reason about, and this amount should be a widely-understood community expectation.

10% seems unreasonable for this amount, and 1% seems like enough to achieve the protocol's goals. Companies are welcome to use a larger percentage if they prefer.

### Why is the tythe based on R&D, not revenue or profit?

 * Modern companies tend to be extremely scalable. Tiny companies can end up generating billions in revenue. Scaling the tithe by revenue or profit ends up not being sensible at either the top or bottom end of the scale.
 * Lots of software-based companies are pre-revenue or profit, but well-funded. It seems reasonable for them to contribute to the maintenance of their open source dependencies.
 * In tech companies, R&D is basically software engineers. It's easy to think about the relationship between R&D budget and the value that open source software provides.
 * R&D is already publicly reported by many companies. This helps with transparency.

### What prevents companies from cheating?

Their public commitment. Companies are collections of people, many of whom are software developers who would want to support this system. This combined with tools that make tythe status visible internally will ensure quite high compliance among companies that make a public commitment.
