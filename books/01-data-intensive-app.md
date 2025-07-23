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