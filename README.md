[![Build Status](https://travis-ci.com/tythe-protocol/go-tythe.svg?branch=master)](https://travis-ci.com/tythe-protocol/go-tythe)

# About

This is a simple CLI client implementing the [Tythe Protocol](https://github.com/aboodman/tythe). It relies on Coinbase for transfer (and eventually, for credit card processing).

# Build

```
git clone https://github.com/tythe-protocol/go-tythe
cd go-tythe
GO111MODULE=on go build ./cmd/go-tythe
```

# Setup

1. Create an account at [Coinbase Pro](https://pro.coinbase.com) if you don't already have one
2. [Create an API key](https://support.pro.coinbase.com/customer/en/portal/articles/2945320-how-do-i-create-an-api-key-for-coinbase-pro-) (all permissions are required)
3. Set the environment variables:
  * `TYTHE_COINBASE_API_KEY`
  * `TYTHE_COINBASE_API_SECRET`
  * `TYTHE_COINBASE_API_PASSPHRASE`
4. Deposit some USDC into your Coinbase Pro account

# Run

```
# This splits a tythe of $2 among all the dependencies of <repo-url>
# For example if you run it with "https://github.com/tythe-protocol/z_test1", it sends $1 to me :-|.
./go-tythe pay-all 2 <repo-url>
```

# Status

Not very much is implemented yet. See [notes](https://github.com/tythe-protocol/go-tythe/blob/master/notes.md) for the plan from here.

# For Open Source Developers

Want to participate in Tythe? Add an appropriate [`tythe.json`](https://github.com/aboodman/tythe/blob/master/tythe-sample.json) to your repo and tweet the URL with `#tythe-protocol` and I'll send you some money.
