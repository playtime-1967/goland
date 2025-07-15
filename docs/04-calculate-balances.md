**“How we calculate balances”**

### ⚙️ Core Concepts & Motivation

1. **Ledger as the source of truth**

   * Monzo uses a **double-entry bookkeeping** system where every transaction is an `EntrySet`—a set of debits and credits—recorded by `service.ledger`.
2. **Address abstraction**

   * An **address** in the ledger represents a money location, defined by: namespace, name, legal entity, currency, account ID.
3. **Balance types**

   * **Customer-facing balance** (available balance seen in app)
   * **Interest-chargeable balance**, based on `committed` timestamps to reflect settled transactions only
4. **Time axes**

   * Used two timestamps per entry: `committed`, `reporting` (when transaction has accounting impact), plus `flake` for ordering

---

### 🕰️ Historical Implementation Challenges

5. **Inconsistent querying**

   * Calling `/balances` with different time axes (`committed` vs `reporting`) led to mismatches
6. **Lack of standard balance definitions**

   * Teams independently queried addresses/time axes, so calculations varied across systems
7. **Poor reconciliation**

   * Data warehouse, backend, and services often produced inconsistent results due to differing filters and time axes
8. **Exposed implementation details**

   * Clients had to know about `flake` IDs, rather than abstracted endpoints

---

### 🔄 The New Balance Architecture

9. **Balance definitions introduction**

   * `service.ledger` now hosts **hardcoded definitions** mapping balance names (e.g., `customer-facing-balance`) to:

     * List of relevant ledger addresses
     * Time axis (`committed` or `reporting`)
10. **Reusing address config**

    * Leverages existing ledger address config, linking definitions to actual addresses
11. **Statically generated BalanceDefinition**

    * Combines name, time axis, and addresses into a generated file used at runtime

---

### 🧮 Calculating a Balance — Workflow

12. **Lookup by balance name**

    * Client requests a named balance (e.g., `customer-facing-balance`)
13. **Fetch EntrySets in parallel**

    * System retrieves entries for all mapped addresses and filters by time axis
14. **Sum all entries**

    * Totals are calculated into a single numeric balance

---

### ✅ Achievements & Improvements

15. **Accuracy**

    * Tests compare against the legacy method to ensure reliability
16. **Consistency**

    * Strong reconciliation between backend and data warehouse
17. **Simplified interface**

    * Clients request with balance name instead of manual address/time axis specs
18. **Abstraction and encapsulation**

    * Internal logic hidden, preventing misuse

---

### 🧠 Why It Matters

19. **Regulatory compliance**

    * Consistent, auditable calculations across systems important for transparency
20. **Developer clarity**

    * Clear contract: “give me balance name, get consistent result” approach avoids ambiguity

---

### 🧩 Summary Flow

| Stage      | Description                                                                                |
| ---------- | ------------------------------------------------------------------------------------------ |
| **Old**    | Clients called `/balances` with custom addresses/time axis resulting in inconsistencies    |
| **New**    | Use named balance definitions → fetch addresses/time axis → query EntrySets → sum → return |
| **Result** | Reliable, auditable, and consistent across systems with simpler client APIs                |
