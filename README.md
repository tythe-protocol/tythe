# Tithe

*Support the church of open source*

## Problem

Open source software is now critical infrastructure for most businesses. Yet most OSS is not actively maintained. Bugs go unaddressed. Critical features go missing and are hacked around. Malicious changes can even get introduced without being noticed.

While the OSS ecosystem is spectacularly creative and prolific, it does not do a good job of incentivizing *maintenance*: the day-to-day work required to keep old software working, adapt it to new targets, update its dependencies, and add new features.

## Proposal

`tithe.json` is a standard way for an open source project to suggest a donation.

*Tithes* are:

* Applicable to for-profit users. Non-profit users aren't expected to pay tithes.
* Enforced socially. There are no license changes. There is no legal enforcement.
* Scaled with revenue. Large businesses pay more than small businesses.
* Open. Money flows directly from users to developers with nobody in the middle.

## Format

```
{
    "destination": {
        # Other supported formats - bitcoin, wire, ach
        "type": "paypal",
        "detail": "maintainer@project.org",
    },
    "base_price_monthly": {
        "currency": "EUR",
        "amount": "2",
    },
}
```

## Calculating the Scaled Tithe

The *scaled tithe* a company should pay is calculated like so:

```
st = base_price_monthly * clamp(R, 1, 1000)
```

Where `R` is the the company's revenue last year, in millions USD. Companies with revenues reported in other currencies should calculate by first converting to equivalent USD at exchange rate from Jan 1 this year.

## Committments

We envision an array of badges that users of open source projects can display on their websites declaring their committment to open source tithes.

A company would display a badge that advertises their promise to support 100% (or more likely 95% or 90%) of their required tithe.

Given the large number of open source developers inside almost all companies, we expect voluntary compliance with these public committments to be quite high.

# Getting Started

There's nothing here now, but I am imaging a variety of software to make calculating, paying, and otherwise working with tithes easier:

1. The very first thing could just be a Patreon integration - users of a software project can click a button on Github to pay via Patreon.
2. A command-line project to find all tithes in a directory and generate a report.
3. A way to do (2) but also make the correct payments.
4. A way to do (2) but show a visual report of where money is going.
