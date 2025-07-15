**“Tolerating full cloud outages with Monzo Stand‑in”**

1. Monzo built **Stand‑in**, a fully independent secondary banking platform, to handle **complete Primary Platform (AWS) outages** ([Monzo][1]).
2. Stand‑in is always running on **Google Cloud Platform (GCP)**, separate from the AWS Primary Platform ([Monzo][1]).
3. It supports essential functionality: **card spending**, **cash withdrawals**, **bank transfers**, **balance checks**, **transaction viewing**, **card freeze/unfreeze** ([Monzo][1]).
4. Each platform has its own **Kubernetes cluster**, **databases**, **queues**, and **locking systems**, with no shared services ([Monzo][1]).
5. Monzo intentionally writes **different implementations** for common behaviors (e.g., payment processing) so Stand‑in won't share code with the Primary, avoiding shared failure vectors ([Monzo][1]).
6. Using fully independent code and infra prevents simultaneous failure during bugs or misconfigurations-–an **isolation against systemic issues** .
7. Stand‑in is cost-effective — running at \~**1% of Primary Platform cost**—because it only supports a **minimal, critical subset** of banking features ([Monzo][1]).
8. Data syncing from Primary to Stand‑in happens via an **event-based “Stand-in Data Syncer”**, which consumes transaction events and applies updates to cards, balances, etc. ([Monzo][1]).
9. All synced data is **immutable**, and Stand‑in tolerates slight lag, which is monitored and alerts are raised if thresholds are exceeded .
10. **Sensitive tokenized data** (e.g., card PANs) is synced via encrypted pipelines, using separate key sets for each platform ([Monzo][1]).
11. When Stand‑in is live, it writes changes (like payments) to an internal queue as **“Monzo Advices”**, not directly to Primary ([Monzo][1]).
12. After stand‑in use, the Primary Platform consumes these Advice events and applies them **verbatim**, ensuring correct final state in the canonical ledger .
13. Stand‑in must merge transactions synced from Primary with ones it generated; this is handled via **correlation IDs** to avoid double-counting ([Monzo][1]).
14. The **Monzo App polls both platforms** and automatically switches UI to a simplified version when Stand‑in is active, exposing only available functionality ([Monzo][1]).
15. Engineers use a **configuration service and CLI** to **enable/disable Stand‑in**, supporting gradual rollbacks and controlled user routing ([Monzo][1]).
16. For **payment routing**, two modes exist:

* **Indirect mode**: Primary proxies requests to Stand‑in (fast, controlled rollout)
* **Direct mode**: Stand‑in processes payments directly (used in full AWS outages) ([Monzo][1], [Monzo Community][2]).

17. Both modes are **tested continuously in production**, ensuring readiness ([Monzo][1]).
18. Stand‑in was used in a **real AWS outage in August 2024**, providing essential banking functionality to all customers while Primary was down for \~1 hour ([Finextra Research][3]).
19. Monzo sees Stand‑in as more than disaster recovery—it's a **resilient, regulatory-forward**, continuously run platform meeting DORA compliance ([Monzo][1]).
20. Future articles will dive into technical details like payment processing, config systems, and testing strategies .

---

### How it works (Client/Server perspective)

* **Clients (Monzo App)** poll both platforms for status, then switch UI accordingly. ✅
* **Primary Platform** emits real-time events continuously; they are consumed by Stand‑in’s Data Syncer. ⬆️
* **Stand‑in Data Syncer** applies those changes into Stand‑in datastore. 🔄
* **User actions on Stand‑in** generate Advices, queued for Primary to reconcile. 📨
* **Primary Platform** processes queued Advices, updates canonical systems, and resyncs resulting transactions back to Stand‑in. 🔃

This architecture ensures continuity, accuracy, and reconciled banking state after an outage — making Monzo uniquely resilient in the event of a full-cloud failure.