# Status

**This project is on indefinite pause while I work on other things. I hope to come back to it someday.**

I think the beginning of the crawl system is pretty cool though. Feel free to lift it.

# Overview

Tythe is a simple, open, and transparent way for companies to fund the Open Source Software they rely on.

Tythe has no middlemen, no chokepoints, no fees, no contracts, and no opaque division of proceeds.

It's an Open Source-style solution to the Open Source maintenance problem.

# How it Works

### 1. Maintainers add a [`.donate`](https://github.com/aboodman/dot-donate) file to the root of their projects

This declares that the projects wants to receive donations (and how those donations should be sent) in a machine-readable format.

### 2. Companies decide an amount to commit to open source

We recommend ["Up to 1%* of R&D"](./covenant.md), but it's your choice.

### 3. Open Source gets funded

Companies run Tythe continuously - either as part of their build process or using [tythe.dev](http://tythe.dev). Tythe monitors dependency trees for participating projects and automatically distributes funds to them. Companies can divide funds however they like, but Tythe provide some reasonable defaults and easy ways to configure.

### CLI

* You can use the Tythe CLI to pay with Bitcoin or USDC
* PayPal support coming soon

### [tythe.dev](http://tythe.dev)

Very early. All you can do so far is show your dependencies. Payment coming soon.

# Setup

First, [Download the latest release](../../releases).

Next, setup payments:

### Coinbase

1. Create an account at [Coinbase Pro](https://pro.coinbase.com) if you don't already have one
2. [Create an API key](https://support.pro.coinbase.com/customer/en/portal/articles/2945320-how-do-i-create-an-api-key-for-coinbase-pro-) (all permissions are required)
3. Set the environment variables:
  * `TYTHE_COINBASE_API_KEY`
  * `TYTHE_COINBASE_API_SECRET`
  * `TYTHE_COINBASE_API_PASSPHRASE`
4. Deposit some USDC and BTC into your Coinbase Pro account

### PayPal

(todo)


# Run

```
# This splits a tythe of $2 among all the dependencies of <repo-url>
# For example if you run it with "https://github.com/tythe-protocol/z_test1", it sends $2 to me.
./tythe pay-all 3 <repo-url>
```

# More

* [Hacking on Tythe](HACKING.md)
* [The Tythe Covenant](covenant.md)
* [About](about.md)
* [FAQ](faq.md)
