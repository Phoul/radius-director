# Architecture

## Overview

RADIUS Director is a configuration platform for FreeRADIUS.

It provides a declarative domain model from which complete, validated FreeRADIUS configurations can be generated.

The architecture is built around a simple principle:

> **Global Objects define reusable resources. Tenants compose those resources into independent RADIUS deployments.**

FreeRADIUS remains the runtime responsible for processing RADIUS requests.

RADIUS Director is responsible for modelling, validating, and generating the configuration.

---

# Architectural Layers

The architecture consists of four logical layers.

```
          Global Objects
                 │
                 ▼
          Tenant Objects
                 │
                 ▼
    Configuration Generator
                 │
                 ▼
      Generated Configuration
                 │
                 ▼
         FreeRADIUS Runtime
```

Each layer has a single responsibility.

---

# Global Objects

Global Objects are reusable definitions that exist once globally.

They may be referenced by one or more tenants.

Examples include:

- Credential Profiles
- Authentication Profiles
- Accounting Profiles
- Monitoring Profiles
- NAS Devices

Global Objects describe reusable behaviour or reusable infrastructure.

They should not contain tenant-specific configuration.

---

# Tenant Objects

A tenant represents a complete, independent RADIUS deployment.

Each tenant owns the infrastructure required for its deployment.

Examples include:

- Database
- RADIUS Servers
- NAS Assignments

Tenant Objects are isolated from one another.

A tenant references Global Objects rather than duplicating them.

---

# Relationship Objects

Some objects exist primarily to describe relationships between other objects.

These objects allow reusable definitions to be combined into tenant-specific configurations.

The primary Relationship Object is:

- NAS Assignment

A NAS Assignment connects:

- a NAS Device
- a Credential Profile
- an Authentication Profile
- an Accounting Profile
- a Monitoring Profile

within the context of a specific tenant.

This allows multiple tenants to reference the same physical NAS while applying different operational policies.

---

# Object Ownership

```
Global Objects
├── Credential Profiles
├── Authentication Profiles
├── Accounting Profiles
├── Monitoring Profiles
└── NAS Devices

Tenants
├── Database
├── RADIUS Servers
└── NAS Assignments
```

Global Objects are shared.

Tenant Objects belong exclusively to a single tenant.

Relationship Objects allow tenants to compose Global Objects without duplication.

---

# Example

```
Global Objects

Credential Profile
Authentication Profile
Accounting Profile
Monitoring Profile
NAS Device

            │
            ▼

Tenant

Database

RADIUS Servers

NAS Assignment
```

The NAS Assignment references the reusable Global Objects while remaining owned by the tenant.

---

# Configuration Generator

The generator transforms the domain model into standard FreeRADIUS configuration.

Generation should always be:

- deterministic
- repeatable
- reproducible
- validated

Generated configuration should never require manual modification.

Changes should always be made to the source model.

---

# Generated Configuration

Generated configuration consists of standard FreeRADIUS configuration files.

Examples include:

- clients.conf
- proxy.conf
- mods-enabled/
- mods-config/
- sites-enabled/

These files are implementation artifacts.

They are not the authoritative source of configuration.

---

# FreeRADIUS Runtime

FreeRADIUS remains responsible for:

- Authentication
- Authorization
- Accounting
- Proxying
- Policy execution

RADIUS Director does not replace FreeRADIUS.

It manages the configuration that FreeRADIUS executes.

---

# Source of Truth

The declarative domain model maintained by RADIUS Director is the authoritative representation of the desired system.

```
Configuration
      │
      ▼
Validation
      │
      ▼
Generation
      │
      ▼
FreeRADIUS
```

Generated configuration should never be edited manually.

All changes should originate from the domain model.

---

# Design Principles

The architecture follows several fundamental principles.

- Standards over Products
- Declarative Configuration
- Infrastructure as Code
- Global Reusable Objects
- Independent Tenant Infrastructure
- Relationship-Based Composition
- Validation Before Generation
- Generated Configuration
- Vendor Neutrality
- Human-Readable Output

These principles are further documented in the project's Architectural Decision Records (ADRs).