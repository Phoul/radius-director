# ADR-0006: Operational Accounting Maintenance

## Status

Accepted

## Context

Some accounting policies require periodic operational activity that cannot be expressed solely as FreeRADIUS request-processing configuration.

One such policy is stale-session cleanup.

RADIUS accounting sessions may remain open when an Accounting-Stop packet is not received.

Authentication-time session verification may determine that the corresponding subscriber session is no longer active and permit the subscriber to reconnect, but this does not necessarily close the stale accounting record.

Leaving stale accounting records open can result in inaccurate accounting state and inaccurate historical association between subscribers and assigned IP addresses.

RADIUS Director therefore requires a mechanism for periodically identifying and closing stale accounting sessions.

Stale-session policy is defined by Accounting Profiles.

Different NAS Assignments within the same tenant may reference Accounting Profiles with different stale-session timeout values.

## Decision

RADIUS Director SHALL support operational accounting maintenance independently of FreeRADIUS request processing.

Operational accounting policy SHALL be defined by the declarative domain model.

The deployment layer SHALL be responsible for provisioning and scheduling the mechanism used to execute operational accounting maintenance.

Stale-session maintenance SHALL operate independently for each tenant.

Each NAS Assignment MAY have its own stale-session timeout through its referenced Accounting Profile.

NAS Assignments whose Accounting Profiles do not define a stale-session timeout SHALL NOT participate in automatic stale-session cleanup.

For the standard FreeRADIUS SQL accounting implementation, stale-session maintenance SHALL use the accounting record's last known accounting update time to determine whether an open session is stale.

A session SHALL be eligible for stale-session cleanup when:

- the session is open
- the session belongs to a NAS Assignment for which stale-session cleanup is enabled
- the last known accounting activity is older than the stale-session timeout defined by that NAS Assignment's Accounting Profile

When stale-session maintenance closes a session, the recorded stop time SHALL represent the last known accounting activity rather than the time at which the maintenance process executes.

Stale-session maintenance SHALL preserve accounting information already recorded from the NAS unless modification is required to mark the session closed.

Maintenance operations SHALL be safe to execute repeatedly.

The initial deployment implementation SHALL schedule accounting maintenance every five minutes.

The maintenance executable itself SHALL perform a single maintenance pass and exit.

Scheduling SHALL remain a deployment responsibility rather than being embedded in the maintenance policy or domain model.

For containerized deployments, operational accounting maintenance SHOULD execute separately from the FreeRADIUS runtime container.

Each tenant SHOULD have an independently scoped maintenance execution environment with access only to the infrastructure required for that tenant.

## Consequences

Advantages:

- stale accounting sessions are automatically closed
- accounting history more accurately reflects the last known subscriber activity
- stale-session cleanup remains independent of authentication-time session verification
- different NAS Devices may use different stale-session timeout policies
- maintenance remains isolated between tenants
- FreeRADIUS runtime containers remain focused on FreeRADIUS
- maintenance execution can be adapted to different deployment platforms
- the same one-shot maintenance executable can be invoked by different scheduling mechanisms
- maintenance can be tested and executed manually without requiring a scheduler

Trade-offs:

- deployments require an additional operational maintenance component
- containerized deployments require a mechanism for periodic execution
- the maintenance component requires access to the tenant's accounting database
- database-specific maintenance implementations may be required for different supported database engines

These trade-offs are preferable to embedding scheduled maintenance into FreeRADIUS configuration or relying on host-specific manual maintenance.

## Deployment Model

Operational accounting maintenance is scoped independently for each tenant.

Example:

```text
Tenant
├── FreeRADIUS Runtime
│
├── Accounting Maintenance
│     └── executes one maintenance pass
│
└── Database
```

The maintenance executable does not determine its own execution schedule.

The deployment implementation determines how it is invoked periodically.

Examples include:

- a scheduler within a dedicated maintenance container
- a Kubernetes CronJob
- a systemd timer
- another platform-native scheduled execution mechanism

The initial maintenance interval is five minutes.

## Stale-Session Semantics

For the standard FreeRADIUS SQL accounting schema, the accounting update timestamp represents the last known accounting activity for an open session.

Conceptually:

```text
open session
    AND
last accounting activity < current time - NAS stale-session timeout
    ↓
session is stale
```

When the stale session is closed:

```text
stop time = last known accounting activity
```

The maintenance process SHALL NOT use its own execution time as the session stop time.

Authentication-time session verification and operational stale-session cleanup solve different problems and SHALL remain independent.

Authentication-time session verification determines whether an apparently active session is actually still present on the NAS.

Operational stale-session cleanup maintains accurate accounting records after accounting activity has ceased.