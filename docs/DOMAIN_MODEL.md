# Domain Model

This document defines the core objects managed by RADIUS Director.

The domain model intentionally describes operational concepts rather than FreeRADIUS configuration files.

Configuration generation is derived from this model.

---

# Platform

Represents a RADIUS Director installation.

A platform may contain one or more tenants.

---

# Tenant

A tenant represents an independent RADIUS deployment.

A tenant owns:

- Credential Profiles
- Authentication Profiles
- Accounting Profiles
- Monitoring Profiles
- NAS definitions
- Database definitions
- RADIUS Servers
- Deployments

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

Represents an external database used by FreeRADIUS.

Examples include:

- Authentication database
- Accounting database

---

# RADIUS Server

Represents a FreeRADIUS instance.

A server consumes generated configuration and participates in one or more deployments.

---

# Deployment

Represents a generated configuration deployed to one or more RADIUS servers.

Future versions may include deployment history and rollback support.