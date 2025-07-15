# Double-entry Bookkeeping Ledger
called append-only logs
accounts, entries, and transactions
ensuring debits = credits in each transaction.

CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL CHECK (type IN ('asset', 'liability', 'equity', 'income', 'expense')),
    created_at TIMESTAMP DEFAULT now()
);

# transactions – Group of entries (one complete event)
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    description TEXT,
    reference TEXT, -- optional external ID or invoice
    posted_at TIMESTAMP NOT NULL DEFAULT now(),
    created_at TIMESTAMP DEFAULT now()
);

# entries – The individual debit or credit lines per transaction
Each transaction must have at least one debit and one credit entry, and the total debits = total credits.
CREATE TABLE entries (
    id SERIAL PRIMARY KEY,
    transaction_id INTEGER NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    account_id INTEGER NOT NULL REFERENCES accounts(id),
    amount NUMERIC(20, 4) NOT NULL CHECK (amount >= 0),
    direction TEXT NOT NULL CHECK (direction IN ('debit', 'credit'))
);


# ✅ Real-world Practice: Balance Calculation

1- Maintain a separate account_balances table.
When a transaction is processed: Update the account_balances table accordingly
CREATE TABLE account_balances (
    account_id INTEGER PRIMARY KEY REFERENCES accounts(id),
    balance NUMERIC(20, 4) NOT NULL DEFAULT 0,
    last_updated_at TIMESTAMP NOT NULL DEFAULT now()
);

2- How to handle reconciliation?
Periodically (e.g. daily, weekly), a background process can:
Recalculate balances from entries for audit/reconciliation, Compare with account_balances, and Raise alert on mismatch.

# So:
Cached balances + Immutable ledger entries is the industry norm.

---
# Overdraft charges 
are fees imposed by banks when a customer spends more money than they have in their account, resulting in a negative balance, or "overdraft". These fees can vary by bank and account type, but can be around $35 per transaction. Overdrafts are essentially a form of short-term borrowing, and banks may also charge interest on the borrowed amount, especially for unarranged overdrafts.


# "effective date" or "ledger date" in entries
posted_at:	When the entry was created in the system
effective_date:	When the entry should take effect

Why they are Useful?
1- Historical Balance Snapshots: WHERE account_id = 1 AND effective_date <= '2025-06-01';
2- Time-Based Reporting: Revenue earned in May but posted in June still counts as May revenue


