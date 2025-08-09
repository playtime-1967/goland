# The Cassandra incident

1. On **July 29, \~13:10 BST**, Monzo experienced major disruption—customers couldn’t log in, view balances or transactions, make payments, withdraw cash, or use in-app support  
2. The incident coincided with a planned **Cassandra cluster scale-up**, adding **six new nodes** to handle growing traffic and capacity needs .
3. Monzo's Cassandra setup comprises **21 servers** with data **replicated three-way**; reads and writes use **quorum consistency**  
4. Pre-incident metrics appeared normal—there were no obvious latency spikes or error surges—so engineers didn’t immediately suspect Cassandra .
5. At **13:14**, internal alerts flagged **Mastercard transaction failures**, and Customer Ops reported issues accessing support tools  
6. Team presumed a bug in the Mastercard service and quickly deployed a **fix at 13:39**, which partially restored card functionality .
7. Shortly after, Monzo’s **internal edge** service began returning **404s for internal endpoints**, hinting at broader access problems  
8. At **14:00**, a deeper investigation revealed that a configuration key was **missing from Cassandra**, so reads returned nonexistent data  
9. Unexpectedly, queries for some data were being served by the **newly added nodes**, which hadn’t yet received their data via streaming  
10. The root cause was traced to `auto_bootstrap = false`: new nodes were **active immediately**, claiming responsibility for partitions without holding data  
11. Monzo had deliberately disabled `auto_bootstrap` to speed up disaster recovery, but this also meant **node activation without data** .
12. As a result, reads to some partitions returned **404 “not found” errors**, causing app failures and transaction inconsistencies .
13. The team began a **manual rollback** by **decommissioning new nodes** one by one, starting \~**14:18**, each taking \~8–10 minutes  
14. By **15:08**, with all new nodes removed, core functionality—including card payments and support—was restored  
15. From **15:08 until 23:00**, data reconciliation occurred via event replay and consistency checking to backfill missing entries .
16. Monzo realized their understanding of `auto_bootstrap` had been flawed: testing scaled to one node, but not six, which broke quorum logic  
17. They’ve now **re-enabled `auto_bootstrap`**, enriched documentation, and updated operational runbooks for Cassandra management  
18. To detect similar failures fast, Monzo plans to surface new metrics like “row not found” counts and alert on unexpected increases  
19. They also plan to **shard the monolithic Cassandra cluster** into smaller, service-specific clusters to limit the blast radius of future changes  
20. The incident taught a clear lesson: **intensive production testing + deep understanding of system flags** are essential before large-scale operations  

---

**In essence**: Monzo triggered a data outage by scaling out Cassandra with misconfigured `auto_bootstrap=false`, causing new nodes to serve empty data partitions. This led to 404s for critical data until the nodes were gracefully decommissioned and data eventually backfilled.