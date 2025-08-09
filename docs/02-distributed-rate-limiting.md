Co‑operative Distributed Rate Limiting

1. Monzo operates \~3,000 microservices on Kubernetes, each typically scaled with multiple replicas  
2. To safely throttle traffic, they require rate limits “per service,” not per replica — because adding replicas multiplies total throughput if each replica applies the same limit .
3. They distinguish two cardinality patterns: **high-cardinality, low-throughput** (e.g., per IP), and **low-cardinality, high-throughput** (e.g., queue consumers)  
4. They focus on the latter pattern: thousands of events/second, few unique limiters  
5. Classic centralized solutions (e.g., Redis counters) were rejected due to network latency on hot paths  
6. Monzo sought a cooperative model where clients enforce limits locally, after obtaining allowance from a central allocator  
7. They drew inspiration from Doorman—a distributed rate limiter that allocates tokens, trusting clients to honor them  
8. They built `service.distrate`, a central service that tracks all interested clients and allocates a share of capacity (`EvenShare` algorithm)  
9. Clients periodically **poll** `service.distrate`, announcing themselves and requesting capacity for configured limiters  
10. The response includes a **lease of capacity** (e.g., 25 events/sec for 5 minutes) and a recommended refresh interval (e.g., 10 seconds)  
11. Clients use this lease to feed a local, in-memory limiter (Golang’s `rate` token bucket), which runs on the hot path  
12. The central service never exceeds total allowed rate and adapts lease sizes as more clients join or leave .
13. `service.distrate` keeps state entirely in memory and uses distributed locks to elect a leader among replicas  
14. If a client fails to refresh before lease expiry, it falls back to a safe “minimal capacity” instead of hard stopping .
15. The naive `EvenShare` algorithm is a starting point; they plan to enhance it (e.g., ramp-up control, customer growth scaling, unused capacity redistribution)  
16. They avoid coupling to Kubernetes API for service discovery—client self-registration suffices  
17. The cooperative model yields **low-latency**, client-driven enforcement without network hops on each request .
18. Server logic is centralized; client logic is minimal and local—this architecture balances correctness and performance .
19. Over time, they intend to support ramp-up smoothing, dynamic allocation based on usage patterns, and algorithm variation .
20. Monzo describes this as a foundational platform they enjoy working with, minimizing opt-outs from auto-scaling and enabling safe horizontal scaling .
