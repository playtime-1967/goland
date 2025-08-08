# Building a Queue on Top of Kafka – Monzo’s Approach

1. **Background** – Monzo has used **Kafka** as a message queue for years, powering over **400 services** across **2,500+ topics**. While Kafka’s *distributed log* architecture is durable and reliable, it’s not naturally suited for *message-oriented queue* semantics where messages are processed independently and acknowledged individually.
2. **Motivation** – The switch from **NSQ** (which lacked durability) to Kafka was driven by the need for **at-least-once delivery guarantees**. However, Kafka’s strict log ordering and partition constraints make concurrent, independent message processing difficult.
3. **Core Solution** – Monzo built a **rich client library** with abstractions that:

   * Support **deadlettering** to skip unprocessable messages.
   * Allow **unordered concurrent consumption** within partitions to improve throughput.
   * Hide Kafka’s low-level complexities from engineers.
4. **Dropping strict ordering** – Most Monzo services don’t require ordering, so relaxing this constraint let them:

   * Prevent a single message from blocking a partition.
   * Enable higher throughput without adding more partitions.
5. **Deadlettering for poison pills** –

   * **Problem:** A “poison pill” message that can’t be processed stalls the entire partition.
   * **Solution:** After `maxAttempts` retries, the message is moved to a **deadletter topic** (specific to the consumer group + topic).
   * Engineers are alerted, and processing continues.
   * For replay, the message is sent to a **retry topic** for that consumer group, avoiding accidental delivery to unrelated consumers.
   * CLI tooling allows engineers to **retry** or **drop** deadlettered messages.
6. **Scaling within partitions** –

   * **Problem:** Partition count is a key scaling lever in Kafka; too few partitions can cause backlog and poor throughput.
   * **Solution:** Introduced **unordered concurrency** so multiple goroutines can process events in the same partition concurrently.
   * Built an **UnorderedOffsetManager** to track in-flight messages and commit the lowest contiguous offset safely.
   * Introduced **max runaway** limit to prevent unbounded backlog when earlier messages block progress.
7. **Impact of unordered concurrency** –

   * Reduces the need for large partition counts (simplifying cluster ops).
   * Gives engineers more **scaling control**: they can simply increase concurrency to handle spikes/backlogs.
8. **Developer experience** –

   * Provided a **typed, fluent API** for consumer configuration:

     ```go
     subscription := balanceproto.BalanceUpdatedTopic.
       NewSubscription("service.account").
       WithDeadletter().
       WithUnorderedConcurrency()
     ```
   * The API enforces **valid combinations at compile time** (e.g., unordered concurrency requires deadlettering).
   * Easily extensible for new options like **distributed locks** or **key-based ordering** within unordered processing.
9. **Ongoing evolution** –

   * Monzo continues to add features to the library.
   * The philosophy remains: **opinionated defaults** and making the “road most travelled” safe and easy for engineers.

---

**In essence:**
Monzo turned Kafka from a strict log into a flexible, message-oriented queue by **sacrificing ordering** for **deadlettering** and **unordered concurrency**. This prevents poison pills from blocking partitions, allows scaling within partitions, and gives engineers more operational control without changing partition counts. The design is reinforced by strong tooling and a compile-time safe API, making it both robust and easy to adopt across hundreds of services.
