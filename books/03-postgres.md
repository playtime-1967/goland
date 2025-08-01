POSTGRESQL


#  "heap" in PostgreSQL
The heap is the physical storage area for table rows.
When you create a table in PostgreSQL, it's stored in an on-disk structure called a heap file.
Unlike some other databases (like clustered indexes in SQL Server or MySQL InnoDB), PostgreSQL does not store table rows in index order — they’re unordered in the heap.

✅ When an index covers all requested data:
PostgreSQL can return results directly from the index.
This is called an index-only scan.
Heap is not touched, so it's faster.

❌ When requested columns aren't in the index:
PostgreSQL uses the index to find row pointers (TIDs).
Then it must go to the heap to fetch the full row.
This is called a heap fetch (or index scan + heap fetch).


CREATE INDEX idx_users_email
ON users (email)
INCLUDE (age, status, city);

# Clustered indexes
PostgreSQL only has non-clustered indexes by default
📌 SQL Server:
Has both:
Clustered indexes: rows in the table are physically sorted by the index key.
Non-clustered indexes: store pointers to the actual rows.
🐘 PostgreSQL:
Does not support true clustered indexes — the table (heap) is always unordered.
Indexes are always separate from the table.

index with included columns:
in addition to storing the full row on the heap, This allows some queries to be answered by using the index alone, without having to look in the heap file. This can make some queries faster, but the duplication of data means the index uses more disk space and slows down writes.


# Consider UUIDs as Promary key if:
You need global uniqueness across systems
You're sharding or syncing across DBs
You want to hide row counts or guessable IDs (for APIs)

