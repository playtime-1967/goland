# NSQ 
realtime distributed messaging
In-memory by default, optional disk
Delivery Semantics:	At-least-once
Use case focus: 
Real-time messaging
Low-latency pipelines (e.g., logs, alerts, metrics)
Systems that don’t require message history/durability

NSQ Limitations:
No message history / replay (unlike Kafka)
Durability is basic — not a long-term log
No transactional guarantees
No built-in ordering


---
CV stories:
high pressure env of banking system, and high pace env of startups!
rabbitmq clustering- move to in-memory in case of incident.
social graph
rust lunching
mongo mig zero-downtime
