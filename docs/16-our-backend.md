# Building a Modern Bank Backend

1. **Why Monzo started with microservices** — From day one Monzo chose a distributed microservices architecture rather than a single monolith so teams could ship independently, scale by area of responsibility, and avoid big-bang releases or long maintenance windows. The goal: 24×7 availability, rapid feature velocity, and the ability to operate in many markets without central coordination.   

2. **Pragmatism for a tiny team** — With only a few backend devs early on, Monzo deliberately picked technologies that moved them fast while keeping “future us” options open. They accepted trade-offs up front (practicality over purity) knowing they could evolve internals later as scale and experience demanded.   

3. **Cluster management & containers — why scheduling mattered** — Monzo moved to a scheduler to treat servers as fungible: run many services per host, reschedule on failure, and elastically scale. Containers (Docker) were used to package apps so the scheduler could manage them uniformly. This led to real operational wins: resilient rescheduling, easier failure testing, and big cost savings (they reported \~75% lower infra costs after switching to Kubernetes).   

4. **Polyglot services, not shared libraries** — Instead of forcing one language or scattering duplicated libraries, Monzo encapsulated cross-cutting concerns as services (e.g., a locking service fronting etcd). This keeps per-language clients small and avoids the maintenance burden of multi-language library forks, while preserving reuse via RPC rather than code-sharing.   

5. **RPC needs and the local proxy (linkerd/Finagle model)** — Monzo needed more than “HTTP only”: routing, retries, smart load balancing, and connection pooling were critical. They adopted an out-of-process proxy approach (linkerd, built on Finagle) run as a local daemon. Apps talk to localhost and the proxy handles retries, advanced load balancing, staged rollouts, and routing logic—so language choice stayed free while ops and resilience improved.   

6. **Asynchronous messaging requirements — why a durable log mattered** — Many tasks (enrichment, notifications, feed insertion) are asynchronous but must never be silently lost. Monzo required: high availability for publishers, horizontal scalability, persistence (so work survives node failures), the ability to replay streams from history, and at-least-once delivery semantics. These shaped their queue choice.   

7. **Why Kafka** — Kafka’s replicated, partitioned commit-log design fit those requirements: durability, replayability, scalable pub/sub (cursors instead of per-consumer queues), and operational robustness for persistent streams. Monzo saw Kafka as the right primitive for long-lived, replayable message streams that underpin financial workflows.   

8. **Operational philosophy: design for recovery & observability** — Monzo emphasised building systems where failures can be handled gracefully: replaying messages, restarting consumers, and operating without human intervention wherever possible. Observability and tooling (to replay, inspect, and recover) are treated as first-class parts of the platform.   

9. **Evolving the platform intentionally** — The post emphasises iterating the platform as needs mature — holding workshops to re-evaluate choices, accepting that some early decisions would be revisited, and documenting learnings to guide future architecture. They flagged follow-up topics (data storage, infra security, hybrid connectivity) as areas they’d detail later.   