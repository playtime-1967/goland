Designing-data-intensive-applications/
There are no solutions, there are only trade-offs.

many applications need to:
(databases)
(caches)
(search indexes)
(stream processing)
(batch processing)

you will learn to ask the right questions to evaluate and compare data systems

# Analytical versus Operational Systems
online transaction processing (OLTP) vs  online analytic processing (OLAP)


# From data warehouse to data lake
A data warehouse often uses a relational data model that is queried through SQL, but it is less well suited to the needs of data scientists, who might need to perform tasks such as:
Transform data into a form that is suitable for training a machine learning model; often this requires turning the rows and columns of a database table into a vector or matrix of numerical values called features. 

# Systems of Record and Derived Data
-Systems of Record: 
also known as source of truth
-Derived data systems:
the result of taking some existing data from another system and transforming or processing it in some way. If you lose derived data, you can recreate it from the original source. A classic example is a cache. Denormalized values, indexes, materialized views, transformed data representations, and models trained on a dataset also fall into this category.

# Distributed Systems
A system that involves several machines communicating via a network
Problems:
Possibility of network failure
Making a call to another service is still vastly slower than calling a function in the same process 
Troubleshooting a distributed system is often difficult
When each service has its own database, maintaining consistency of data across those different services becomes the application’s problem.

# Serverless, or function-as-a-service (FaaS)
the cloud provider automatically allocates and frees hardware resources as needed, based on the incoming requests to your service
----

# Chapter 2. Defining Nonfunctional Requirements

#Case Study: Social Network Home Timelines
1- Let’s assume that users make 500 million posts per day, or 5,700 posts per second on average. Occasionally, the rate can spike as high as 150,000 posts/second [4]. Let’s also assume that the average user follows 200 people and has 200 followers.

2- Imagine we keep all of the data in a relational database as shown in Figure 2-1. We have one table for users, one table for posts, and one table for follow relationships.
SELECT posts.*, users.* FROM posts
  JOIN follows ON posts.sender_id = follows.followee_id
  JOIN users   ON posts.sender_id = users.id
  WHERE follows.follower_id = current_user
  ORDER BY posts.timestamp DESC
  LIMIT 1000
3- let’s assume that after somebody makes a post, we want their followers to be able to see it within 5 seconds. One way of doing that would be for the user’s client to repeat the query above every 5 seconds while the user is online (this is known as polling). If we assume that 10 million users are online and logged in at the same time, that would mean running the query 2 million times per second. Even if you increase the polling interval, this is a lot.

4- the query above is quite expensive: if you are following 200 people, it needs to fetch a list of recent posts by each of those 200 people, and merge those lists. 2 million timeline queries per second then means that the database needs to look up the recent posts from some sender 400 million times per second


#How can we do better?
Firstly, instead of polling, it would be better if the server actively pushed new posts to any followers who are currently online.
Secondly, we should precompute the results of the query above so that a user’s request for their home timeline can be served from a cache.

Imagine that for each user we store a data structure containing their home timeline, i.e., the recent posts by people they are following. Every time a user makes a post, we look up all of their followers, and insert that post into the home timeline of each follower.

The downside of this approach is that we now need to do more work every time a user makes a post, because the home timelines are derived data that needs to be updated. 
At a rate of 5,700 posts posted per second, if the average post reaches 200 followers, we will need to do just over 1 million home timeline writes per second. This is a lot, but it’s still a significant saving compared to the 400 million per-sender post lookups per second.

If the rate of posts spikes due to some special event, we don’t have to do the timeline deliveries immediately—we can enqueue them and accept that it will temporarily take a bit longer for posts to show up in followers’ timelines. Even during such load spikes, timelines remain fast to load, since we simply serve them from a cache.
This process of precomputing and updating the results of a query is called materialization, and the timeline cache is an example of a materialized view.

The cost of writes for most users is modest, but a social network also has to consider some extreme cases:
When a celebrity account with a very large number of followers makes a post, we have to do a large amount of work to insert that post into the home timelines of each of their millions of followers. In this case it’s not okay to drop some of those writes. One way of solving this problem is to handle celebrity posts separately from everyone else’s posts:

# 1- Describing Performance
two main types of metric:
- Response time: The elapsed time from the moment when a user makes a request until they receive the answer.
- Throughput: The number of requests per second, or the data volume per second, that the system is processing. The unit of measurement is “somethings per second/micro/etc”.

There is often a connection between throughput and response time; The service has a low response time when request throughput is low, but response time increases as load increases. This is because of queueing: when a request arrives on a highly loaded system, it’s likely that the CPU is already in the process of handling an earlier request, and therefore the incoming request needs to wait until the earlier request has been completed.


# When an overloaded system won’t recover
If a system is close to overload, with throughput pushed close to the limit, it can sometimes enter a vicious cycle where it becomes less efficient and hence even more overloaded. For example, if there is a long queue of requests waiting to be handled, response times may increase so much that clients time out and resend their request. This causes the rate of requests to increase even further, making the problem worse—a retry storm. Even when the load is reduced again, such a system may remain in an overloaded state until it is rebooted or otherwise reset. This phenomenon is called a metastable failure, and it can cause serious outages in production systems [7, 8].

To avoid retries overloading a service, you can increase and randomize the time between successive retries on the client side (exponential backoff [9, 10]), and temporarily stop sending requests to a service that has returned errors or timed out recently (using a circuit breaker [11, 12] or token bucket algorithm [13]). The server can also detect when it is approaching overload and start proactively rejecting requests (load shedding [14]), and send back responses asking clients to slow down (backpressure [1, 15]). The choice of queueing and load-balancing algorithms can also make a difference [16].

throughput determines the required computing resources (e.g., how many servers you need), and hence the cost of serving a particular workload. If throughput is likely to increase beyond what the current hardware can handle, the capacity needs to be expanded; a system is said to be scalable if its maximum throughput can be significantly increased by adding computing resources.


# Latency and Response Time
- The response time is what the client sees; it includes all delays incurred anywhere in the system.
- The service time is the duration for which the service is actively processing the user request.
- Queueing delays can occur at several points in the flow: for example, after a request is received, it might need to wait until a CPU is available before it can be processed; a response packet might need to be buffered before it is sent over the network if other tasks on the same machine are sending a lot of data via the outbound network interface.
- Latency is a catch-all term for time during which a request is not being actively processed, i.e., during which it is latent. In particular, network latency or network delay refers to the time that request and response spend traveling through the network.


The response time can vary significantly from one request to the next, even if you keep making the same request over and over again. Many factors can add random delays: for example, a context switch to a background process, the loss of a network packet and TCP retransmission, a garbage collection pause, a page fault forcing a read from disk, mechanical vibrations in the server rack

# average response time  vs percentiles
- average response timeis not a very good metric if you want to know your “typical” response time, because it doesn’t tell you how many users actually experienced that delay.
- use percentiles: take your list of response times and sort it from fastest to slowest, then the median is the halfway point

#High percentiles of response times, also known as tail latencies:
Amazon describes response time requirements for internal services in terms of the 99.9th percentile, even though it only affects 1 in 1,000 requests. This is because the customers with the slowest requests are often those who have the most data on their accounts because they have made many purchases—that is, they’re the most valuable customers


# 2- Reliability and Fault Tolerance
means continuing to work correctly, even when things go wrong;
To be more precise about things going wrong, we will distinguish between faults and failures;
Fault:
A fault is when a particular part of a system stops working correctly: for example, if a single hard drive malfunctions, or a single machine crashes, or an external service (that the system depends on) has an outage.
Failure:
A failure is when the system as a whole stops providing the required service to the user;


# Fault Tolerance
We call a system fault-tolerant if it continues providing the required service to the user in spite of certain faults occurring. If a system cannot tolerate a certain part becoming faulty, we call that part a single point of failure (SPOF), because a fault in that part escalates to cause the failure of the whole system.

#Exactly-once semantics
in the context of messaging and data processing, refers to the guarantee that a message or operation is processed exactly one time, without duplication or loss, even in the presence of failures

#Fault injection
increase the rate of faults by triggering them deliberately
Although we generally prefer tolerating faults over preventing faults, there are cases where prevention is better than cure (e.g., because no cure exists). This is the case with security matters, for example: if an attacker has compromised a system and gained access to sensitive data, that event cannot be undone.

culture of blameless postmortems:
after an incident, the people involved are encouraged to share full details about what happened, without fear of punishment, since this allows others in the organization to learn how to prevent similar problems in the future

# 3- Scalability
a system’s ability to cope with increased load
Even if a system is working reliably today, that doesn’t mean it will necessarily work reliably in the future.

#Case Study: a Startup, that currently only has a small number of users
- it is counterproductive to worry about hypothetical scale that might be needed in the future: in the best case, investments in scalability are wasted effort and premature optimization; in the worst case, they lock you into an inflexible design and make it harder to evolve your application.
- scalability is not a one-dimensional label: it is meaningless to say “X is scalable”, instead considering questions like:
“If the system grows in a particular way, what are our options for coping with the growth?”
“How can we add computing resources to handle the additional load?”
“Based on current growth projections, when will we hit the limits of our current architecture?”

# Describing Load
Often this will be a measure of throughput, for example, the number of requests per second to a service, how many gigabytes of new data arrive per day, or the number of shopping cart checkouts per hour. Sometimes you care about the peak of some variable quantity, such as the number of simultaneously online users.
ratio of reads to writes in a database, the hit rate on a cache

linear scalability
If you can double the resources in order to handle twice the load, while keeping performance the same, and this is considered a good thing.
Much more likely is that the cost grows faster than linearly. For example, if you have a lot of data, then processing a single write request may involve more work than if you have a small amount of data, even if the size of the request is the same.


# Shared-Nothing Architecture
also called horizontal scaling or scaling out
The advantages of shared-nothing are that it has the potential to scale linearly, it can more easily adjust its hardware resources as load increases or decreases, and it can achieve greater fault tolerance by distributing the system across multiple data centers and regions. The downsides are that it requires explicit sharding.

# Principles for Scalability
a system that is designed to handle 100,000 requests per second, each 1 kB in size, looks very different from a system that is designed for 3 requests per minute, each 2 GB in size—even though the two systems have the same data throughput (100 MB/sec).
- A good general principle for scalability is to break a system down into smaller components that can operate largely independently from each other. 

# 4- Maintainability
It is widely recognized that the majority of the cost of software is not in its initial development, but in its ongoing maintenance— fixing bugs, keeping its systems operational, investigating failures, adapting it to new platforms, modifying it for new use cases, repaying technical debt, and adding new features.

#Operability: Making Life Easy for Operations

#Simplicity: Managing Complexity
break it down into two categories, essential and accidental complexity.
-essential complexity: is inherent in the problem domain of the application
-accidental complexity arises only because of limitations of our tooling.

One of the best tools we have for managing complexity is abstraction.
For example, high-level programming languages are abstractions that hide machine code, CPU registers, and syscalls.
SQL is an abstraction that hides complex on-disk and in-memory data structures, concurrent requests from other clients, and inconsistencies after crashes.

#Evolvability: Making Change Easy
---
# Chapter 3. Data Models and Query Languages

The limits of my language mean the limits of my world.
Ludwig Wittgenstein

Data models are perhaps the most important part of developing software: not only on how the software is written, but also on how we think about the problem that we are solving.

#one-to-few rather than one-to-many relationship (Nosql):
- a résumé typically has a small number of positions. 
- comments on a celebrity’s social media post, of which there could be many thousands—embedding them all in the same document may be a lot, so the relational approach is preferable.

#normalization
- The downside of a normalized representation is that every time you want to display a record containing an ID, you have to do an additional lookup to resolve the ID into something human-readable.
- normalized data is usually faster to write (since there is only one copy), but slower to query. denormalized data is usually faster to read (fewer joins), but more expensive to write (more copies to update, more disk space used).
- in systems of small to moderate scale, a normalized data model is often best,

hydrating the IDs
looking up the human-readable information by ID is called hydrating the IDs

#Denormalization in the social networking case study
- the implementation of materialized timelines at X does not store the actual text of each post. each entry actually only stores the post ID, the ID of the user who posted it, and a little bit of extra information to identify reposts and replies.
- The reason for storing only IDs in the precomputed timeline is that the data they refer to is fast-changing: the number of likes and replies may change multiple times per second on a popular post, and some users regularly change their username or profile photo.
- This example shows that having to perform joins when reading data is not, as sometimes claimed, an impediment to creating high-performance, scalable services.
- the social network case study shows that the choice is not immediately obvious: the most scalable approach may involve denormalizing some things and leaving other things normalized.

associative table
also called join table- many-to-many relationships
Many-to-many relationships often need to be queried in “both directions” :for example, finding all of the organizations that a particular person has worked for, and finding all of the people who have worked at a particular organization.

The document model limitations: 
- you cannot refer directly to a nested item within a document. If you do need to reference nested items, a relational approach works better, since you can refer to any item directly by its ID.

- schemaless is misleading. a more accurate term is schema-on-read (the structure of the data is implicit, and only interpreted when the data is read). So it's schemaless-on-write.

#Data locality for reads and writes
it is generally recommended that you keep documents fairly small and avoid frequent small updates to a document.

#Graph-Like Data Models
one-to-many relationships (tree-structured data) 
many-to-many relationships(graph)
A graph consists of two kinds of objects: vertices (also known as nodes or entities) and edges (also known as relationships or arcs):
-Social graphs
Vertices are people, and edges indicate which people know each other.
-The web graph
Vertices are web pages, and edges indicate HTML links to other pages.
-Road or rail networks
Vertices are junctions, and edges represent the roads or railway lines between them.

Well-known algorithms:
map navigation apps search for the shortest path between two points in a road network.

- each vertex consists of:
A unique identifier
A label (string) to describe what type of object this vertex represents
A set of outgoing edges
A set of incoming edges
A collection of properties (key-value pairs)

- Each edge consists of:
A unique identifier
The vertex at which the edge starts (the tail vertex)
The vertex at which the edge ends (the head vertex)
A label to describe the kind of relationship between the two vertices
A collection of properties (key-value pairs)

The Cypher Query Language
is a query language for property graphs

#GraphQL
allow clients to request a JSON document with a particular structure, containing the fields necessary for rendering its user interface.  

#Event Sourcing and CQRS

---
# Chapter 4. Storage and Retrieval

two families of storage engines for OLTP
- log-structured storage engines that write out immutable data files
- storage engines such as B-trees that update data in-place. 

log:
an append-only sequence of records on disk. It doesn’t have to be human-readable; it might be binary and intended only for internal use by the database system.
the cost of a lookup is O(n)
index
In order to efficiently find the value for a particular key in the database, we need a different data structure: an index.
the general idea is to structure the data in a particular way (e.g., sorted by some key) that makes it faster to locate the data you want.
Any kind of index consumes additional disk space and usually slows down writes, because the index also needs to be updated every time data is written.

hash table
Range queries are not efficient. For example, you cannot easily scan over all keys between 10000 and 19999—you’d have to look up each key individually in the hash map.

Sorted String Table, SSTable 
stores key-value pairs, but it ensures that they are sorted by key, and each key only appears once in the file.

SSTable with a sparse index
you do not need to keep all the keys in memory: you can group the key-value pairs within an SSTable into blocks of a few kilobytes, and then store the first key of each block in the index. This kind of index, which stores only some of the keys, is called sparse

Memtable
an LSM-tree (Log-Structured Merge-Tree) structure
the Memtable serves as an in-memory write-back cache for recent write operations


checksum 
used to ensure data integrity, meaning they help verify that data is consistent and hasn't been changed accidentally or maliciously.

B-Trees
The log-structured approach is popular, but it is not the only form of key-value storage. The most widely used structure for reading and writing database records by key is the B-tree.
Like SSTables, B-trees keep key-value pairs sorted by key, which allows efficient key-value lookups and range queries.

The log-structured indexes break the database down into variable-size segments, typically several megabytes or more in size, that are written once and are then immutable. By contrast, B-trees break the database down into fixed-size blocks or pages, and may overwrite a page in-place. A page is traditionally 4 KiB in size

# Comparing B-Trees and LSM-Trees
As a rule of thumb, LSM-trees are better suited for write-heavy applications, whereas B-trees are faster for reads.

backpressure in LSM-Trees
suspend all reads and writes until the memtable has been written out to disk

Write in LSM-trees
a value is first written to the log for durability, then again when the memtable is written to disk, and again every time the key-value pair is part of a compaction

Disk space usage
B-trees can become fragmented over time: for example, if a large number of keys are deleted, the database file may contain a lot of pages that are no longer used by the B-tree. Subsequent additions to the B-tree can use those free pages, but they can’t easily be returned to the operating system because they are in the middle of the file, so they still take up space on the filesystem. Databases therefore need a background process that moves pages around to place them better, such as the vacuum process in PostgreSQL.

Fragmentation is less of a problem in LSM-trees, since the compaction process periodically rewrites the data files anyway, and SSTables don’t have pages with unused space.

secondary index 
the indexed values are not necessarily unique; that is, there might be many rows under the same index entry. 

Storing values within the index
If the actual data (row, document, vertex) is stored directly within the index structure, it is called a clustered index.

in-memory storages
Even a disk-based storage engine may never need to read from disk if you have enough memory, because the operating system caches recently used disk blocks in memory anyway. Rather, they can be faster because they can avoid the overheads of encoding in-memory data structures in a form that can be written to disk.

----
# Chapter 5. Encoding and Evolution
staged rollout(rolling upgrade)
deploying the new version to a few nodes at a time, checking whether the new version is running smoothly, and gradually working your way through all the nodes. This allows new versions to be deployed without service downtime, and thus encourages more frequent releases and better evolvability.

need to maintain compatibility in both directions@
-Backward compatibility
Newer code can read data that was written by older code.
-Forward compatibility
Older code can read data that was written by newer code.

# Encoding
Programs usually work with data in (at least) two different representations:
- In memory, data is kept in objects, structs, lists, arrays, hash tables, trees, and so on. These data structures are optimized for efficient access and manipulation by the CPU (typically using pointers).

- When you want to write data to a file or send it over the network, you have to encode it as some kind of self-contained sequence of bytes (for example, a JSON document).
Thus, we need some kind of translation between the two representations. The translation from the in-memory representation to a byte sequence is called encoding(also known as serialization or marshalling) and the reverse is called decoding (parsing, deserialization, unmarshalling).

There are exceptions in which encoding/decoding is not needed:
zero-copy data formats that are designed to be used both at runtime and on disk/on the network, without an explicit conversion step

#Language-Specific Formats
The encoding is often tied to a particular programming language, and reading the data in another language is very difficult. 

When moving to standardized encodings that can be written and read by many programming languages, JSON and XML are the obvious contenders.
There is a lot of ambiguity around the encoding of numbers. In XML and CSV, you cannot distinguish between a number and a string. JSON distinguishes strings and numbers, but it doesn’t distinguish integers and floating-point numbers, and it doesn’t specify a precision.
JSON and XML have good support for Unicode character strings (i.e., human-readable text), but they don’t support binary strings (sequences of bytes without a character encoding).

#Binary encoding
JSON is less verbose than XML, but both still use a lot of space compared to binary formats. 

#Protocol Buffers
(protobuf) is a binary encoding library developed at Google
Protocol Buffers requires a schema for any data that is encoded. 
Protocol Buffers comes with a code generation tool that takes a schema definition and produces classes that implement the schema in various programming languages.
Each field is identified by its tag number (the numbers 1, 2, 3 in the sample schema)

# Modes of Dataflow
Via databases
Via service calls: REST and RPC
Via workflow engines: Durable Execution
Via asynchronous messages

The problems with remote procedure calls (RPCs)
A local function call is predictable and either succeeds or fails.
A network request is unpredictable
The client and the service may be implemented in different programming languages, so the RPC framework must translate datatypes from one language into another.

#Durable Execution and Workflows

# Chapter 6. Replication

Replication 
means keeping a copy of the same data on multiple machines that are connected via a network
If the data that you’re replicating does not change over time, then replication is easy: you just need to copy the data to every node once, and you’re done. 
All of the difficulty in replication lies in handling changes to replicated data

Backups and replication
replicas quickly reflect writes from one node on other nodes, but backups store old snapshots of the data so that you can go back in time. If you accidentally delete some data, replication doesn’t help since the deletion will have also been propagated to the replicas

Each node that stores a copy of the database is called a replica. With multiple replicas, a question inevitably arises: how do we ensure that all the data ends up on all the replicas?


# Single-Leader Replication
 Whenever the leader writes new data to its local storage, it also sends the data change to all of its followers as part of a replication log or change stream
 When a client wants to read from the database, it can query either the leader or any of the followers.

Synchronous Versus Asynchronous Replication
the replication to follower 1 is synchronous: the leader waits until follower 1 has confirmed that it received the write before reporting success to the user, and before making the write visible to other clients.
The replication to follower 2 is asynchronous: the leader sends the message, but doesn’t wait for a response from the follower.


# Multi-Leader Replication
Imagine you have a database with replicas in several different regions (perhaps so that you can tolerate the failure of an entire region, or perhaps in order to be closer to your users). This is known as a geographically distributed, geo-distributed or geo-replicated setup.

every write can be processed in the local region and is replicated asynchronously to the other regions. Thus, the inter-region network delay is hidden from users, which means the perceived performance may be better.


Sync Engines and Local-First Software
an application that needs to continue to work while it is disconnected from the internet.
If you make any changes while you are offline, they need to be synced with a server and your other devices when the device is next online.
In this case, every device has a local database replica that acts as a leader (it accepts write requests), and there is an asynchronous multi-leader replication process (sync) between the replicas of your calendar on all of your devices.
From an architectural point of view, this setup is very similar to multi-leader replication between regions, taken to the extreme: each device is a “region,” and the network connection between them is extremely unreliable.

Real-time collaboration, offline-first, and local-first apps
This again results in a multi-leader architecture: each web browser tab that has opened the shared file is a replica, and any updates that you make to the file are asynchronously replicated to the devices of the other users who have opened the same file.
If multiple users have changed the file concurrently, conflict resolution logic may be needed to merge those changes.
Git is a local-first collaboration system

Pros of sync engines
when using a sync engine, you have persistent state on the client, and communication with the server is moved into a background process.
A sync engine combined with a reactive programming model is a good way of implementing this 
Sync engines work best when all the data that the user may need is downloaded in advance and stored persistently on the client.

cons of sync engines
The biggest problem with multi-leader replication—both in a geo-distributed server-side database and a local-first sync engine on end user devices—is that concurrent writes on different leaders can lead to conflicts that need to be resolved.

#Dealing with Conflicting Writes
concurrent writes on different leaders can lead to conflicts that need to be resolved.
- Conflict avoidance
One strategy for conflicts is to avoid them occurring in the first place. 
for ex: 
If you have two leaders, you could set them up so that one leader only generates odd numbers and the other only generates even numbers.
-Last write wins (discarding concurrent writes)
If conflicts can’t be avoided, the simplest way of resolving them is to attach a timestamp to each write, and to always use the value with the greatest timestamp... This achieves the goal that eventually all replicas end up in a consistent state, but at the cost of data loss.
a problem with LWW is that if a real-time clock (e.g. a Unix timestamp) is used as timestamp for the writes, the system becomes very sensitive to clock synchronization. If one node has a clock that is ahead of the others, and you try to overwrite a value written by that node, your write may be ignored as it may have a lower timestamp. This problem can be solved by using a logical clock.
If lost updates are not acceptable, you need to use one of the conflict resolution approaches:
-Manual conflict resolution


# Leaderless Replication
like Cassandra
this kind of database is also known as Dynamo-style.
In some leaderless implementations, the client directly sends its writes to several replicas, while in others, a coordinator node does this on behalf of the client. However, unlike a leader database, that coordinator does not enforce a particular ordering of writes. 

Catching up on missed writes
The replication system should ensure that eventually all the data is copied to every replica. After an unavailable node comes back online

Quorums for reading and writing
In our example (n = 3, w = 2, r = 2.) As long as w + r > n
You can think of r and w as the minimum number of votes required for the read or write to be valid.

In Dynamo-style databases, the parameters n, w, and r are typically configurable. A common choice is to make n an odd number (typically 3 or 5) and to set w = r = (n + 1) / 2 (rounded up). However, you can vary the numbers as you see fit. For example, a workload with few writes and many reads may benefit from setting w = n and r = 1. This makes reads faster, but has the disadvantage that just one failed node causes all database writes to fail.

The quorum condition, w + r > n, allows the system to tolerate unavailable nodes as follows:
If w < n, we can still process writes if a node is unavailable.
If r < n, we can still process reads if a node is unavailable.
With n = 3, w = 2, r = 2 we can tolerate one unavailable node.
With n = 5, w = 3, r = 3 we can tolerate two unavailable nodes.


#Single-Leader vs. Leaderless Replication Performance
Reading from the leader ensures up-to-date responses, but it suffers from performance problems:
Read throughput is limited by the leader’s capacity to handle requests 
If the leader fails, you have to wait for the fault to be detected

A big advantage of a leaderless architecture is that it is more resilient against such issues.

#Multi-region operation
Cassandra and ScyllaDB implement their multi-region support within the normal leaderless model


# Chapter 7. Sharding
split up a large amount of data into smaller shards or partitions, and store different shards on different nodes
each piece of data (each record, row, or document) belongs to exactly one shard. 

Sharding is usually combined with replication so that copies of each shard are stored on multiple nodes.

What we call a shard in this chapter has many different names depending on which software you’re using: it’s called a partition in Kafka, a range in CockroachDB, a region in HBase and TiDB, a tablet in Bigtable and YugabyteDB, a vnode in Cassandra.

in PostgreSQL, partitioning is a way of splitting a large table into several files that are stored on the same machine (which has several advantages, such as making it very fast to delete an entire partition), whereas sharding splits a dataset across multiple machines

Pros and Cons of Sharding
The primary reason for sharding a database is scalability: it’s a solution if the volume of data or the WRITE throughput has become too great for a single node to handle, as it allows you to spread that data and those writes across multiple nodes. (If READ throughput is the problem, you don’t necessarily need sharding—you can use READ SCALING.

While replication is useful at both small and large scale, because it enables fault tolerance and offline operation, sharding is a heavyweight solution that is mostly relevant at large scale.

sharding often adds complexity: you typically have to decide which records to put in which shard by choosing a partition key; all records with the same partition key are placed in the same shard.

Another problem with sharding is that a write may need to update related records in several different shards. While transactions on a single node are quite common, ensuring consistency across multiple shards requires a distributed transaction.


Some systems use sharding even on a single machine, typically running one single-threaded process per CPU core to make use of the parallelism in the CPU- like Redis.

Sharding for Multitenancy
Software as a Service (SaaS) products and cloud services are often multitenant, where each tenant is a customer. Multiple users may have logins on the same tenant, but each tenant has a self-contained dataset that is separate from other tenants.

Using sharding for multitenancy has several advantages:
Resource isolation
Permission isolation
Per-tenant backup and restore
Regulatory compliance- GDPR
Data residence

Our goal with sharding is to spread the data and the query load evenly across nodes- every node takes a fair share

1- Sharding by Key Range
Within each shard, keys are stored in sorted order (e.g., in a B-tree or SSTables).
A downside of key range sharding is that you can easily get a hot shard if there are a lot of writes to nearby keys.

Rebalancing key-range sharded data:
When you first set up your database, there are no key ranges to split into shards. Some databases, such as HBase and MongoDB, allow you to configure an initial set of shards on an empty database, which is called pre-splitting. 

An advantage of key-range sharding is that the number of shards adapts to the data volume. If there is only a small amount of data, a small number of shards is sufficient, so overheads are small; if there is a huge amount of data, the size of each individual shard is limited to a configurable maximum.

A downside of this approach is that splitting a shard is an expensive operation, since it requires all of its data to be rewritten into new files. A shard that needs splitting is often also one that is under high load, and the cost of splitting can exacerbate that load, risking it becoming overloaded.

2- Sharding by Hash of Key
first hash the partition key before mapping it to a shard.
The problem with the mod N(Key Range) approach is that if the number of nodes N changes, most of the keys have to be moved from one node to another.
In this model(Hash of Key), only entire shards are moved between nodes, which is cheaper than splitting shards. 

3- Sharding by hash range
combine key-range sharding with a hash function so that each shard contains a range of hash values rather than a range of keys.


# Request Routing
very similar to service discovery

Three Methods:
1- Allow clients to contact any node (e.g., via a round-robin load balancer). If that node coincidentally owns the shard to which the request applies, it can handle the request directly; otherwise, it forwards the request to the appropriate node, receives the reply, and passes the reply along to the client.
2- Send all requests from clients to a routing tier first, which determines the node that should handle each request and forwards it accordingly. This routing tier does not itself handle any requests; it only acts as a shard-aware load balancer.
3- Require that clients be aware of the sharding and the assignment of shards to nodes. In this case, a client can connect directly to the appropriate node, without any intermediary.


consensus algorithms
Many distributed data systems rely on a separate coordination service which They use consensus algorithms:
Many distributed data systems rely on a separate coordination service such as ZooKeeper or etcd to keep track of shard assignments.
Each node registers itself in ZooKeeper, and ZooKeeper maintains the authoritative mapping of shards to nodes. Other actors, such as the routing tier or the sharding-aware client, can subscribe to this information in ZooKeeper. Whenever a shard changes ownership, or a node is added or removed, ZooKeeper notifies the routing tier

# Sharding and Secondary Indexes
