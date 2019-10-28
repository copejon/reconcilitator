# reconcilitator
_because I just can't be bother to find that one missing entry_

Still a WIP.

## Purpose

This project is meant to facilitate reconciling my bank accound to my YNAB budget by comparing one exported CSV against one another and
reporting unmatched entries.  Using this report, I can quickly find and either add or correct the date of entries in the budget.
Besides streamlininbg bugdet reconciliation, it also provides for a mostly<sup>*</sup> accurate parity between each 
register's running balance (if your using [Toolkit for YNAB]( https://www.toolkitforynab.com/), which you should).


> \* Same-day entries are ordered differently for every register, so only running balance at the final entry for a day should 
> be expected to match.  Not the running balance for each entry.


## Principle

Comparing Payee to Payee and the amount isn't a workable solution because the Payee field values change between registers.
For instance, in the bank register Payee: Wal-mart 01011231 Main St. Foovile, Tx, while in YNAB Payee: Walmart.  Instead,
the idea is that each day will have a given number of entries.  So, we can hash our entries into days and then compare amounts.
If the same day in both registers does not have the same number of entries, that's a red flag. If an entry is matched, it's marked.  
This prevents false positives from occuring when a day contains 2 or more charges of the same amount*.

> \* This does mean that psuedo-false negatives could be reported in these cases. e.g.
> 
> BANK|Amt|YNAB|Amt
> ----|----|----|----
> Walm 1231|$100|Walmart|$100
> Target| $100 | n/a | n/a
> 
> The app may match the Bank:Target charge to the YNAB:Walmart charge and the Report the Bank:Walmart charge as missing.
> If you run into this, look for another entry in that day with the same dollar amount.  It's probably the real missing entry.

See [TODO issue](https://github.com/copejon/reconcilitator/issues/2) for a running list of stuff
