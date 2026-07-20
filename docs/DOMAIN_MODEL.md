# Domain Model

This document defines the core objects managed by RADIUS Director.

The domain model intentionally describes operational concepts rather than FreeRADIUS configuration files.

Configuration generation is derived from this model.

---

# Platform

The platform contains reusable Global Objects that may be referenced by one or more tenants.

Global Objects include:

- Credential Profiles
- Authentication Profiles
- Accounting Profiles
- Monitoring Profiles
- NAS Devices

---

# Tenant

Each tenant contains tenant-specific infrastructure.

Tenant Objects include:

- Database
- RADIUS Servers
- NAS Assignments

Tenants reference Global Objects rather than duplicating them.

---

# Relationship Objects

Some objects exist primarily to describe relationships between other objects.

Relationship Objects allow reusable Global Objects to be composed into tenant-specific deployments.

The primary Relationship Object is:

- NAS Assignment

A NAS Assignment references:

- NAS Device
- Credential Profile
- Authentication Profile
- Accounting Profile
- Monitoring Profile

while remaining owned by a single tenant.

---

# Credential Profile

Defines reusable authentication credentials.

Typical properties include:

- RADIUS shared secret
- CoA shared secret

Multiple NAS devices may reference the same credential profile.

---

# Authentication Profile

Defines authentication behaviour.

Examples include:

- Simultaneous Use
- Session verification
- SQL policies
- Vendor-specific behaviour

---

# Accounting Profile

Defines accounting behaviour.

Examples include:

- Accounting storage
- Interim update handling
- Session cleanup
- Retention policies

---

# Monitoring Profile

Defines operational monitoring.

Examples include:

- Ping tests
- Authentication tests
- CoA tests
- SNMP verification

---

# NAS

Represents a RADIUS client.

Typical properties:

- Name
- Address
- Vendor
- Credential Profile
- Authentication Profile
- Monitoring Profile

---

# Database

Represents the primary database backing a tenant's RADIUS deployment.

Typical properties include:

- Database engine
- Host
- Port
- Database name
- Authentication credentials
- TLS configuration

Each tenant owns its own database definition.

Future versions may support additional backend services where appropriate.

---

# RADIUS Server

Represents a FreeRADIUS instance.

A server consumes generated configuration and participates in one or more deployments.

---

# Deployment

Represents a generated configuration deployed to one or more RADIUS servers.

Future versions may include deployment history and rollback support.