## Building a Processing System for International Payments

1. **Launch & Scale**

   * In **April 2023**, Monzo enabled users to receive Euros via IBAN. Since then, they've expanded support to **over 40 additional currencies**. To manage this expansion cleanly, they built a brand-new **International Payment Processing system**, which went live just two weeks before the blog post.
     

2. **What Counts as an International Payment?**

   * Any bank transfer crossing country or currency borders—like someone in Australia sending AUD to a UK GBP account—qualifies as an international payment. Monzo supports these via IBAN and handles any necessary currency conversion.
     

3. **How It Works: Correspondent Banking**

   * Monzo doesn’t hold 40+ foreign currency accounts. Instead they partner with **correspondent banks**, which receive payments on Monzo’s behalf, perform FX conversion, and forward the funds over SWIFT. The correspondent routes the transfer efficiently using SWIFT’s SSIs database.
     

4. **Design Principles Guiding the System**
   Monzo emphasized three core principles in the processor's design:

   * **Correctness**: Every payment must process **once and exactly once**, with strong idempotency, reconciliation, and coherence mechanisms.
   * **Testability**: They built components small and decoupled. Each—adapter, decider, effector—can be tested independently for clarity and confidence.
   * **Scalability & Reusability**: A common, modular architecture lets them onboard new partners or payment types without rewriting the processing engine.
     

5. **Core Processing Stages**
   The system follows a multi-stage flow:

   * **Adaptor**: Translates partner-specific formats (e.g., SWIFT messages) into a **common payment representation**. Handles validation, parsing, and even callbacks.
   * **Decisioning (Deciders)**: Check whether to **Accept**, **Reject**, or **Hold** each payment—e.g., verifying IBAN exists, account is open, compliance passes, no limits are breached. Payment outcomes follow strict precedence rules: Hold > Reject > Accept.
   * **Effect Generation**: Creates the actions for accepted payments—such as updating ledger balance, notifying the sender, emitting events.
   * **Effect Application**: Executes those actions, interacts with downstream services, logs metadata for auditability.
     

6. **Onboarding New Partners or Payment Types**

   * To add a **new partner**, engineers write a new **adaptor service** (for message formats, APIs, etc.), while reusing existing decisioning and effect logic.
   * To roll out a **new payment type**, they reuse deciders and effectors; only configuration adjustments are needed. All platform-level insights—metrics, dashboards, alerts, runbooks—automatically apply.
     

7. **Multi-Tiered Testing Approach**

   * **Unit tests** cover individual adaptors, deciders, effectors, and the control flow among them—ensuring each piece behaves as expected in isolation.
   * **Acceptance tests** validate the entire processing workflow end-to-end by initiating example payments and confirming real-world behavior, especially across service dependencies. These run regularly in staging to catch integration issues early.
     

8. **Engineering Impact & Philosophy**

   * The new processor streamlines adding or changing international payment capabilities across many currencies and service providers.
   * It also improves developer experience by encouraging modularity, reusability, and clear testing—making the system safer and faster to evolve.
     

---

**In essence**: Monzo built a modular, testable, and reusable system to process international payments reliably. Leveraging **correspondent banking**, plus a clear pipeline—adaptors, deciders, effectors, applications—they ensure correctness, ease of integration, and extensibility. Strong testing layers and a consistent architecture ensures they can scale to new currencies and partners with confidence.