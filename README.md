# reconcilitator
_because I just can't be bothered to find that one missing entry_

Still a WIP.

## Purpose

This project is meant to facilitate reconciling my bank accound to my YNAB budget by comparing one exported CSV against one another and
reporting unmatched entries.  Using this report, I can quickly find and either add or correct the date of entries in the budget.
Besides streamlininbg bugdet reconciliation, it also provides for a mostly<sup>*</sup> accurate parity between each 
register's running balance (if you're using [Toolkit for YNAB]( https://www.toolkitforynab.com/), which you should).

> \* Same-day entries are ordered differently for every register, so only running balance at the final entry for a day should 
> be expected to match.  Not the running balance for each entry.

## Principle

This app operates on this assumption: there should always be entry parity by date.  That is, for every day of transactions, the
number of those transactions should be identical between the bank register and your budget's register.  In implementation,
this means that we first map all entries by date, then perform comparisons by amount.

This design was taken becase comparing Payee to Payee and the Transaction amount isn't a trival solution. Payee field values change between registers,
tor instance -- the bank register Payee: Wal-mart 01011231 Main St. Foovile, Tx, while the YNAB Payee: Walmart.   

#### Limitations

The principle essentially is a Best Guess solution where we cannot be 100% certain of matches.  This opens up the possibility
for false positives and false negatives to sneak in.  Fortunately, these only occur under 1 scenario: when a day has more than one entry
with identical transaction amounts _and_ those transactions are not in parity between registers.

E.g.:

BANK|Amt|YNAB|Amt
----|----|----|----
Walma 1231|$100|Walmart|$100
Target| $100 | n/a | n/a

The app may match the Bank:Target charge to the YNAB:Walmart charge and the Report the Bank:Walmart charge as missing.
When looking through the audit output, it should be immediately apparent that the Uncleared charge is actually present.  At this point
find the next charge in that day with the exact same transaction. That will most likely be the actual missing entry.

See [TODO issue](https://github.com/copejon/reconcilitator/issues/2) for a running list of stuff
