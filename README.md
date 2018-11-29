# Tithe

*Support the church of open source*

## Problem

Open source software is now critical infrastructure for most businesses. Yet most OSS is not actively maintained. Bugs go unaddressed. Critical features go missing and are hacked around. Malicious changes can even get introduced without being noticed.

The key issue is that OSS doesn’t incentivize anyone to do the grungy day-to-day work to keep critical software infrastructure working correctly.

## Proposal

`tithe.json` is a standard way for an open source project to suggest a donation.

*Tithes* are:

* Applicable to for-profit users. Non-profit users aren't expected to pay tithes.
* Enforced socially. There are no license changes. There is no legal enforcement.
* Scaled with revenue. Large businesses pay more than small businesses.

## Format

```
{
    "destination": {
        # Other supported formats - bitcoin, wire, ech
        "type": "paypal",
        "detail": "maintainer@project.org",
    },
    "base_price_monthly": {
        "currency": "EUR",
        "amount": "5",
    },
}
```

## Calculating the Scaled Tithe

The *scaled tithe* a company should pay is calculated like so:

```
st = base_price_monthly * clamp(R, 1, 1000)
```

Where `R` is the the company's revenue last year, in millions USD. Companies with revenues reported in other currencies should calculate by first converting to equivalent USD at exchange rate from Jan 1 this year.

## Badges

We envision an array of badges that users of open source projects can display on their websites declaring their level of support for tithes.

A company would commit publicly to support 100% of tithes in their codebase. Or more realistically, 90% or 80%.

There could even be audit services to prove the committment, though I'm not sure that's necessary.

## Implementations

TODO - there would eventually go here links to software packages one could run against a directory tree to find all tithes, compute a report, pay tithes, or even do so continuously.
