Commitlog, Memtable, and SSTable are important components of Cassandra architecture

# Cassandra Memtable Definition
Write operations are first written to the Cassandra Memtable structure rather than directly to disk. This speeds write operations

A Memtable in Cassandra is a data structure that serves as a write-back cache for recent in-memory writes. When write operations occur, data is stored in the Memtable first.
Later, when the Cassandra Memtable data size reaches a certain threshold, its contents are flushed to disk as sorted string tables (SSTables) for persistent storage to optimize write performance.

# Commitlog 
ensures data durability
Commitlog is a recovery mechanism used to ensure data durability in case of node failure or crash

When the memtable gets bigger than some threshold—typically a few megabytes—write it out to disk in sorted order as an SSTable file. 

# UPDATE operation
UPDATE is treated as a new write, not an in-place modification. Updates are effectively insertions with newer timestamps.
-Write goes to the Memtable (an in-memory write-back cache).
-An entry is also appended to the Commit Log (for durability).
-The Memtable holds the latest value — including "updates".
-Eventually, the Memtable is flushed to a new SSTable.

# DELETE 
-It writes a tombstone — a special marker that says: "this row is deleted."
-The tombstone is stored in the Memtable, then flushed to an SSTable.
-During reads, Cassandra sees the tombstone and hides the deleted data.
-If the tombstone is older than gc_grace_seconds (default: 10 days), it's permanently purged.

# Cassandra compaction 
is a background process that reorganizes data stored in SSTables to improve read performance and reduce disk space usage. It involves merging multiple SSTables into a single, more efficient SSTable, discarding outdated data and tombstones, and optimizing data storage. Compaction is crucial for maintaining Cassandra's performance and efficiency. 