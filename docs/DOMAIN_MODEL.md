# Domain Model

This document defines the core objects managed by RADIUS Director.

The domain model intentionally describes operational concepts rather than FreeRADIUS configuration files.

Configuration generation is derived from this model.

---

# Global Objects

Global Objects are reusable objects that may be referenced by one or more tenants.

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

Multiple NAS Assignments may reference the same Credential Profile.

---

# Authentication Profile

Defines authentication behaviour.

Examples include:

- Simultaneous Use
- Session verification
- SQL policies
- Vendor-specific behaviour

Authentication Profiles are referenced by NAS Assignments.

---

# Accounting Profile

Defines accounting behaviour.

Examples include:

- Accounting storage
- Interim update handling
- Session cleanup
- Retention policies

Accounting Profiles are referenced by NAS Assignments.

---

# Monitoring Profile

Defines operational monitoring.

Examples include:

- Ping tests
- Authentication tests
- CoA tests
- SNMP verification

Monitoring Profiles are referenced by NAS Assignments.

---

# NAS Device

Represents a physical or virtual RADIUS client.

A NAS Device is a reusable Global Object that describes the device itself, independent of how it is used by any particular tenant.

Typical properties include:

- Name
- Address
- Vendor
- Model
- Description
- Tags

Operational behaviour such as credentials, authentication, accounting, and monitoring is defined through NAS Assignments rather than directly on the NAS Device.

---

# NAS Assignment

Represents how a tenant uses a NAS Device.

A NAS Assignment is a Relationship Object that combines reusable Global Objects into a tenant-specific configuration.

Each NAS Assignment references:

- NAS Device
- Credential Profile
- Authentication Profile
- Accounting Profile
- Monitoring Profile

Multiple tenants may reference the same NAS Device while applying different operational policies through separate NAS Assignments.

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

Represents a FreeRADIUS instance belonging to a tenant.

A RADIUS Server consumes generated configuration and services authentication, authorization, accounting, proxying, and Change of Authorization (CoA) requests for the tenant.

Typical properties include:

- Name
- Hostname or Address
- Role
- Description