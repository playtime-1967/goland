## Vertically Scaling Ordered Consumption Using Kafka

1. **Problem with Horizontal Scaling**
   While increasing Kafka topic partitions is a common scaling strategy, Monzo identified several drawbacks:

   * **Operational complexity**: Adding partitions requires careful ordering and configuration to avoid skipping messages, especially with default consumer offsets set to `offsetNewest`  
   * **Resource overhead**: Partitions consume cluster metadata, CPU, and complicate rebalances  
   * **Irreversibility**: Adding partitions is a one-way operation; you can’t easily remove them later  
   * **Mismatch with deployment preferences**: Monzo favors running fewer, larger pods—more partitions would push them toward more, smaller pods  
   * **Static backlog issues remain**: Increasing partitions doesn’t help clear existing processing backlogs quickly  

2. **Extending Unordered Concurrency with Partition-Key Ordering**
   Monzo enhanced its existing unordered concurrent subscription to maintain ordering **within partition keys**, while enabling concurrency **across different keys**  

3. **Design Elements for Per-Key Ordered Concurrency**

   * **Partition key usage**: Ensures all messages for the same key land in the same partition, enabling safe, ordered processing  
   * **Distributed lock per partition**: Prevents multiple consumers from processing a partition simultaneously during rebalances  
   * **Local, per-key locks**: Synchronize access so that only one message per partition key is in flight, preserving in-key ordering without distributed coordination overhead  

4. **Concurrency Behavior in Action**
   In a scenario with four concurrent goroutines:

   * Three process unique keys (e.g., `acc_123`, `acc_456`, `acc_789`).
   * A fourth is blocked from processing a message for `acc_123` until the key is unlocked (i.e., previous message committed).
     This allows concurrent processing across keys with in-order guarantees per key, built atop the unordered concurrency model with “max runaway” safeguards  

5. **Limitations & Rationale**

   * The naive blocking behavior—where a locked key halts all picking of that message—is intentional. Monzo anticipates that key collisions within short time windows are rare  
   * The design leans on simplicity, deferring more complex handling unless real-world use reveals issues  

6. **Dependence on Deadletter Support**
   As with unordered concurrency, this per-key ordering feature requires deadletter support. Monzo prefers halting on rare errors rather than stalling the queue indefinitely, and deadlettering ensures safety and observability  

---

**In essence**:
Monzo’s evolved Kafka strategy allows engineers to scale vertically—with concurrency within fewer partitions—while still maintaining strict ordering for messages sharing the same key. They achieve this through smart use of partition key placement, locks, and coordinated logic, all while retaining the safety features of deadlettering and simplified scaling experience across their systems.