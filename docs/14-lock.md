# Mutual Exclusion
ensure that if a process is already performing write operation on a data object [critical section] no other process/thread is allowed to access/modify the same object until the first process has finished writing upon the data object.
this alg prevents race condition.

# Local Lock:
like tokio::sync::Mutex
Uses OS or language runtime primitives (like semaphores, mutexes, spinlocks).
Only coordinates threads within the same memory space or OS process.
Very fast.

# Distributed Lock:
Across multiple processes and machines in a network or cluster.
Examples:
Databases: Using SELECT ... FOR UPDATE or advisory locks
Key-value stores: Redis Redlock, etcd, ZooKeeper.
Cloud services: AWS DynamoDB conditional writes, Google Cloud Spanner transactions.
How it works:
Uses an external, shared, reliable system that all processes can access.
The "lock" state is stored in that system so all participants see the same lock state.

# Some concepts in distributed systems
Clock Drift
refers to mismatches in timekeeping across different machines (e.g., process A’s clock is a few milliseconds ahead of process B’s).
In lock systems that rely on timeouts or leases (like Redlock in Redis), unsynchronized clocks can cause one client to think a lock has expired when another client is still operating under it—a recipe for inconsistency. 

Failover
the automatic (or manual) switch to a redundant component when the primary one fails.
This ensures the system continues operating even if one node fails—known as maintaining liveness.

Split-Brain
happens when a distributed system is partitioned (e.g., due to network failure), and multiple segments each believe they are the “only active leader.”
this can lead to multiple consumers processing the same partition simultaneously—breaking exclusivity and possibly causing data duplication or inconsistency.




