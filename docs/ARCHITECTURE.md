# Architecture

## Overview

RADIUS Director is a declarative deployment platform for FreeRADIUS.

It models multi-tenant RADIUS infrastructure using a declarative domain model from which complete, version-aware FreeRADIUS deployments are generated.

The architecture is built around a simple principle:

> **Global Objects define reusable resources. Tenants compose those resources into independent FreeRADIUS deployments.**

Each tenant represents a complete, independent FreeRADIUS configuration tree.

The generated configuration for one tenant is entirely isolated from every other tenant and can be deployed independently.

RADIUS Director is responsible for modelling, validating, generating, and deploying complete FreeRADIUS environments.

FreeRADIUS remains responsible for processing RADIUS requests at runtime.

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
      Configuration Model
                 │
                 ▼
      Generation Pipeline
                 │
                 ▼
 Per-Tenant Configuration Trees
                 │
                 ▼
        Deployment Layer
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

A tenant represents a complete, independent FreeRADIUS deployment.

Each tenant owns the infrastructure required for its deployment.

Examples include:

- Database
- RADIUS Server
- NAS Assignments

The RADIUS Server object includes the target FreeRADIUS version.

The configured version determines:

- validation rules
- managed templates
- generated configuration
- deployment artifacts

Tenant Objects are isolated from one another.

A tenant references Global Objects rather than duplicating them.

Each tenant ultimately produces its own complete FreeRADIUS configuration tree.

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
├── RADIUS Server
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

RADIUS Server

NAS Assignment
```

The NAS Assignment references the reusable Global Objects while remaining owned by the tenant.

---

# Configuration Generator

The generator transforms the validated domain model into standard FreeRADIUS configuration.

Generation is performed independently for each tenant.

Each tenant produces a complete, self-contained FreeRADIUS configuration tree.

Generation should always be:

- deterministic
- repeatable
- reproducible
- validated

Generated configuration should never require manual modification.

Changes should always be made to the source model.

---

# Generated Configuration

Each tenant generates its own independent FreeRADIUS configuration tree.

For example:

```
output/
├── customer-a/
│   ├── clients.conf
│   ├── mods-available/
│   ├── mods-enabled/
│   ├── sites-available/
│   └── sites-enabled/
│
└── customer-b/
    ├── clients.conf
    ├── mods-available/
    ├── mods-enabled/
    ├── sites-available/
    └── sites-enabled/
```

Each generated tree can be deployed independently.

Generated configuration consists entirely of standard FreeRADIUS configuration files.

Examples include:

- clients.conf
- proxy.conf
- mods-available/
- mods-enabled/
- mods-config/
- sites-available/
- sites-enabled/

These files are implementation artifacts.

They are not the authoritative source of configuration.

---

# Deployment

Deployment transforms a generated configuration tree into a runnable FreeRADIUS environment.

Responsibilities include:

- selecting the configured FreeRADIUS version
- selecting version-compatible managed templates
- generating deployment artifacts
- provisioning supporting infrastructure
- deploying or updating FreeRADIUS instances

The deployment layer is intentionally separated from configuration generation.

Generation produces deterministic configuration.

Deployment determines how that configuration is executed.

---

# FreeRADIUS Runtime

Each generated configuration tree is intended to be executed by an independent FreeRADIUS instance.

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
Per-Tenant Configuration Trees
      │
      ▼
Deployment
      │
      ▼
Running FreeRADIUS
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
- One Configuration Tree Per Tenant
- Relationship-Based Composition
- Validation Before Generation
- Deterministic Generation
- Generated Configuration
- Vendor Neutrality
- Human-Readable Output
- Version-Aware Deployments
- Separation of Generation and Deployment

These principles are further documented in the project's Architectural Decision Records (ADRs).