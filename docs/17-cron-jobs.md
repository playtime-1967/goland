## a Better System for Scheduling Cron Jobs

1. **What’s a cron job at Monzo?**
   They run scheduled tasks—like processing BACS files or health checks—by hitting service endpoints automatically, just like any other microservice workload. Only difference: activation is time-based, not human-triggered.
    

2. **Early solutions and inefficiencies**

   * **Long-running cron containers**: Always-on pods that periodically make HTTP requests—wasteful when idle most of the time.
   * **Kubernetes CronJobs**: Better—they spin up pods, run the job, and shut down. But still incur significant overhead (pod startup, scheduling), especially as cron count grows.
      

3. **Developer experience friction**
   Monzo prioritizes delightful tooling. But the cron setup was clunky:

   * Required writing cryptic crontab expressions (`"30 * * * 1-5 sh /run.sh"`)
   * Lacked timezone awareness—K8s CronJobs run in UTC, causing unexpected fire times during daylight savings.
   * Opaque deployment—no easy feedback on cron deployment, failure, or status. Engineers had to build observability themselves.
      

4. **Monzo’s “dream cron system” goals**

   * Define crons directly in **Go**, alongside service code.
   * Express schedules simply—e.g., “every X minutes/hours/days”—instead of crontab syntax.
   * Get **Slack alerts** for cron creation, updates, and failures.
   * Support granular **retries**, **timeouts**, **paging**, and error handling.
   * Register crons as part of service deployment, keeping lifecycle in sync.
   * Provide a **delightful CLI** to list, inspect, pause crons—especially handy during incidents.
      

5. **Declarative cron configuration in Go**
   The new library, `cronfig`, allows expressive cron definitions in Go:

   ```go
   cron.Config {
     CronName:    "some-job",
     Description: "Call service foo at 15:30 in Europe/London",
     Request:     fooproto.SomeRequest{},
     Schedule: cron.Schedule {
       Crontab:  "30 15 * * 1-5",
       Timezone: "Europe/London",
     },
   }
   ```

   Or a simpler "every 10 minutes" cron using `OncePerDuration`, which avoids the classic “top-of-the-hour” traffic bursts by spreading executions naturally.
    

6. **Handling “top-of-the-hour” clustering pain**
   Traditional crontabs cause many jobs to fire at the same time (e.g., exactly on the hour). `OncePerDuration` solves this elegantly by using the cron’s creation timestamp as the baseline—creating a natural, randomized spread without manual offsets.
    

7. **Failure semantics built in**
   Cron definitions include user-friendly failure handling, including:

   * `RpcTimeout` (how long to wait for service response),
   * `ExecutionWindow` (how long past schedule is acceptable),
   * `PageOnFailure` (alerting),
   * `RunbookURL` (manual troubleshooting link).
      

8. **Deployment-triggered cron registration**
   Cron definitions in code are **code-generated** into JSON and registered with Monzo’s cron service during deployment:

   * Adds/updates cron schedules if changed.
   * If unchanged, does nothing.
   * Automatically applies defaults or polyfills with every deploy—keeping cron configs up to date implicitly.
      

9. **Execution engine**
   The cron service polls every **60 seconds**, executes due crons, recalculates next run times, and schedules future runs. On failure, it still schedules the next run and sends alerts to cron owners.
    