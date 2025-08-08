## Building a Reactive Fraud Prevention Platform

1. Published **30 July 2025**, Monzo shared how and why they redesigned their Fraud Prevention Platform ([Monzo][1]).

2. Fraud is a massive problem—In 2024, UK Finance estimated fraud losses at **£1.17 billion** ([Monzo][1]).

3. Fraud detection is tricky:

   * Fraudsters are sophisticated and agile—“**whack-a-mole**” style evolving attacks ([Monzo][1]).
   * Only **1 in 10,000 transactions** is fraud, so precision is critical to avoid false positives or missed fraud ([Monzo][1]).

4. Key system requirements:

   * **Scale** with complexity of controls
   * **Rapid deployment** to react to new threats
   * **Visibility** to monitor control performance
   * Must also handle **millions of transactions daily** with **low latency** and **fault tolerance** ([Monzo][1]).

---

## System Design at a Glance

Monzo's system follows **four primary steps** for every transaction:

1. **Select Controls** – Choose the relevant controls, e.g., apply ML model only on bank transfers ([Monzo][1]).
2. **Load Features** – Compute inputs needed for controls (e.g., payment amount) ([Monzo][1]).
3. **Execute Controls** – Run detectors (ML models), actions, and selection logic to decide if intervention is needed ([Monzo][1]).
4. **Apply Actions** – Based on control output, intervene (e.g., block or warn), raising the necessary alerts or user notifications ([Monzo][1]).

---

## Components That Power the Platform

### Engine

* A Go microservice containing a **Controls Repository** and **Executor**.
* Fraud controls are written in **Starlark** as **pure functions**, enabling easy **backtesting** via BigQuery ([Monzo][1]).

### Control Types

* **Detectors** – ML models that signal fraud likelihood.
* **Action Controls** – Decide what to do when fraud is detected.
* **Action Selection Control** – Aggregates and selects final actions in a scalable modular pipeline ([Monzo][1]).

### Feature Loader

* Uses a **DAG** (Directed Acyclic Graph) to compute features:

  * **Just-in-time**: computed on the fly (e.g., payment reference)
  * **Near real-time**: cached/precomputed (e.g., today’s spending)
  * **Batch**: precomputed offline via SQL (e.g., yearly spending) ([Monzo][1]).
* DAG ensures **resilience**—errors in one node don’t break the request, and timeouts stop high-latency computations ([Monzo][1]).

### Action Applier

* Applies interventions based on determined actions.
* Stateful: tracks prior actions to avoid redundant interventions.
* Emits metrics to BigQuery for monitoring and includes **rate limits** to avoid systemic bugs causing mass actions ([Monzo][1]).

---

## In Summary

Monzo’s fraud platform supports **complex, rapidly-changing controls** across high-volume transactions with **low latency**. Its modular design—pure-function controls, DAG-based feature computation, and robust action handling—enables scalable, testable, and resilient fraud prevention.