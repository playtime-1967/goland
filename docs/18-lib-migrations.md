## Run Migrations Across 2,800 Microservices

1. **Scale & Central Ownership**
   Monzo manages a staggering **2,800 microservices**, all written in Go and maintained in a **monorepo**, to ensure consistency in dependencies and tooling. ([monzo.com](https://monzo.com/blog/how-we-run-migrations-across-2800-microservices))
   (, [InfoQ][2])

2. **Core Migration Principles**

   * **Centrally driven**: A single team runs migrations centrally to reduce coordination overhead and prevent stalled or inconsistent rollouts.
   * **Zero downtime**: Since many services are critical, migrations must not interrupt service.
   * **Gradual roll-forward**: Deploy incrementally to minimize blast radius and allow fast rollback if needed.
   * **80/20 automation**: Heavily automate common patterns, manually handle edge cases for efficiency.
     ()

3. **Process & Strategy Overview**

   **a. Planning & Alignment**

   * Migrations begin with proposals shared via Slack, open for feedback.
   * High-risk changes go through synchronous architecture reviews to catch edge cases early.
     ()

   **b. Wrapping the Old Library**

   * They create a wrapper around the existing (old) library to dynamically switch between implementations via configuration, enabling toggling without redeploying all services.
   * This also enables telemetry insertion and offers a more opinionated interface.
     ()

   **c. Refactoring Call Sites**

   * The most commonly referenced functions/types are updated via automated tools (e.g., `gopls`, `gorename`).
   * Rare call sites are tackled manually or replaced with a generic API, keeping the wrapper minimalist.
   * CI checks (e.g., semgrep) block new dependencies on the old library.
     ()

   **d. Adding the New Library**

   * Integrate the new library behind the same wrapper, initially disabled via config to allow smooth rollout.
     ()

   **e. Mass Deployment**

   * Build tooling to asynchronously deploy changes across all services.
   * Services are tagged by criticality; lower-risk services are deployed first.
   * Automated rollback checks are in place to catch regressions early.
     (, [shivangsnewsletter.com][3])

   **f. Configuration-based Rollout Control**

   * Instead of enabling new behavior via code deploy, they use their live config system.
   * Services refresh config every **60 seconds**, allowing quick activation or rollback of the new library per service, per user set, or percentage rollout.
     (, [InfoQ][2])

   **g. Final Cleanup**

   * Once stable and widely adopted, they remove the old library from the wrapper entirely.
     ()

4. **Bonus Insight: Opinionated Wrapper Interface**
   In Reddit, Monzo engineers explain they prefer wrapping third-party APIs to provide a constrained, safer interface rather than exposing the full, complex API surface. This aligns with their migration strategy.
   ([Reddit][4])

---

**In essence:** Monzo handles library migrations across thousands of microservices through centralized planning, automated tooling, gradual rollout, configuration-driven toggles, and wrapped interfaces—enabling safe, consistent, and low-risk evolution across its entire fleet.