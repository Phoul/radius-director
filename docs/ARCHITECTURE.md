# Architecture

## Overview

RADIUS Director is a declarative deployment platform for FreeRADIUS.

It models multi-tenant RADIUS infrastructure using a declarative domain model from which version-aware FreeRADIUS deployments are generated.

RADIUS Director generates and manages only the deployment-specific configuration required to realize the desired deployment. The underlying FreeRADIUS distribution, its default configuration, and runtime components remain the responsibility of FreeRADIUS.

The architecture is built around a simple principle:

> **Global Objects define reusable resources. Tenants compose those resources into independent FreeRADIUS deployments.**

Each tenant represents a complete, independent FreeRADIUS deployment.

The generated configuration for one tenant is entirely isolated from every other tenant and can be deployed independently.

RADIUS Director is responsible for modelling, validating, generating, and deploying managed FreeRADIUS deployments.

FreeRADIUS remains responsible for processing RADIUS requests at runtime.

---

# Architectural Layers

The architecture consists of several logical layers.

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
   Managed Configuration Trees
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
- Trusted RADIUS Clients

Global Objects describe reusable behaviour or reusable infrastructure.

They should not contain tenant-specific configuration.

---

# Tenant Objects

A tenant represents a complete, independent FreeRADIUS deployment.

Each tenant owns the infrastructure required for its deployment.

Examples include:

- Database
- RADIUS Server

Relationship Objects owned by the tenant include:

- NAS Assignments
- Trusted RADIUS Client Assignments

The RADIUS Server object includes the target FreeRADIUS version.

The configured version determines:

- validation rules
- managed templates
- generated configuration
- deployment artifacts

Tenant Objects are isolated from one another.

A tenant references Global Objects rather than duplicating them.

Each tenant ultimately produces its own managed FreeRADIUS configuration tree.

---

# Relationship Objects

Some objects exist primarily to describe relationships between other objects.

These objects allow reusable definitions to be combined into tenant-specific configurations.

Relationship Objects include:

- NAS Assignment
- Trusted RADIUS Client Assignment

A NAS Assignment connects:

- a NAS Device
- a Credential Profile
- an Authentication Profile
- an Accounting Profile
- a Monitoring Profile

within the context of a specific tenant.

A Trusted RADIUS Client Assignment connects:

- a Trusted RADIUS Client
- a Credential Profile

within the context of a specific tenant.

This allows multiple tenants to reference the same Global Objects while applying tenant-specific operational configuration.

---

# Object Ownership

```
Global Objects
├── Credential Profiles
├── Authentication Profiles
├── Accounting Profiles
├── Monitoring Profiles
├── NAS Devices
└── Trusted RADIUS Clients

Tenants
├── Database
├── RADIUS Server
├── NAS Assignments
└── Trusted RADIUS Client Assignments
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

Each tenant produces a complete, self-contained managed FreeRADIUS configuration tree.

Generation should always be:

- deterministic
- repeatable
- reproducible
- validated

Generated configuration should never require manual modification.

Changes should always be made to the source model.

---

# Generated Configuration

Each tenant generates its own independent managed FreeRADIUS configuration tree.

For example:

```
output/
├── customer-a/
│   ├── clients.conf
│   ├── radiusd.conf
│   ├── proxy.conf
│   ├── mods-available/
│   │   └── sql
│   ├── mods-enabled/
│   │   └── sql -> ../mods-available/sql
│   ├── sites-available/
│   │   ├── default
│   │   └── coa
│   └── sites-enabled/
│       ├── default -> ../sites-available/default
│       └── coa -> ../sites-available/coa
└── customer-b/
    ├── clients.conf
    ├── radiusd.conf
    ├── proxy.conf
    ├── mods-available/
    │   └── sql
    ├── mods-enabled/
    │   └── sql -> ../mods-available/sql
    ├── sites-available/
    │   ├── default
    │   └── coa
    └── sites-enabled/
        ├── default -> ../sites-available/default
        └── coa -> ../sites-available/coa
```

Each generated tree can be deployed independently.

Generated configuration consists of the managed subset of standard FreeRADIUS configuration files required for the deployment.

Examples include:

- clients.conf
- radiusd.conf
- proxy.conf
- mods-available/
- sites-available/

These files are implementation artifacts.

They are not the authoritative source of configuration.

---

# Managed Configuration

RADIUS Director manages only the configuration files required to describe a deployment.

It does not replace or regenerate the complete FreeRADIUS configuration distributed with a FreeRADIUS installation.

Typical managed files include:

- clients.conf
- radiusd.conf
- proxy.conf
- mods-available/
- sites-available/

The managed configuration tree may also include symbolic links required to enable managed modules or virtual servers, including:

- mods-enabled/
- sites-enabled/

The following remain the responsibility of the installed FreeRADIUS distribution:

- dictionaries
- certificates
- policy libraries
- mods-config/
- documentation
- other runtime resources

This separation minimizes maintenance, preserves compatibility with upstream FreeRADIUS releases, and allows RADIUS Director to focus on deployment-specific configuration.

---

# Ownership Boundaries

```
Declarative Model
        │
        ▼
RADIUS Director
        │
        ▼
Managed Configuration
        │
        ▼
FreeRADIUS Distribution
        │
        ▼
Running FreeRADIUS
```

The declarative model is the source of truth.

RADIUS Director transforms the declarative model into managed configuration.

The installed FreeRADIUS distribution provides the runtime, default configuration, and supporting resources.

Each layer has a clearly defined owner and responsibility.

---

# Deployment

Deployment transforms a generated configuration tree into a runnable FreeRADIUS environment.

Deployment responsibilities include:

- installing the managed configuration tree
- provisioning supporting infrastructure
- provisioning operational maintenance mechanisms
- starting or updating FreeRADIUS instances

The deployment layer is intentionally separated from configuration generation.

Generation produces deterministic configuration.

Deployment determines how that configuration is executed.

---

# Operational Maintenance

Some domain policies require periodic operational activity rather than FreeRADIUS configuration generation.

RADIUS Director distinguishes between:

- the policy that defines the desired operational behaviour
- the mechanism used to execute that behaviour

Operational policy is defined by the declarative domain model.

The deployment layer is responsible for provisioning the mechanism required to execute that policy.

For example, an Accounting Profile may define a `stale_session_timeout`.

The timeout determines when an accounting session associated with a NAS Assignment is considered stale.

Different NAS Assignments within the same tenant may reference Accounting Profiles with different stale-session timeout values.

RADIUS Director uses the NAS identity associated with accounting records to apply the stale-session policy belonging to the corresponding NAS Assignment.

The mechanism used to periodically evaluate and close stale sessions is an operational maintenance responsibility rather than FreeRADIUS request-processing configuration.

The deployment layer is responsible for arranging periodic execution of accounting maintenance appropriate to the deployment environment.

Examples may include:

- a system service and timer
- a scheduled container
- a platform-native scheduled job

The scheduling mechanism is not part of the domain model.

Operational maintenance should be deterministic and safe to execute repeatedly.

Stale-session maintenance operates independently of authentication-time session verification.

Authentication-time session verification may determine whether a subscriber is currently active and permit reconnection when an existing accounting record is stale.

Accounting maintenance is responsible for ensuring that stale accounting records are eventually closed and that accounting history remains internally consistent.

When stale-session maintenance closes an accounting session, the recorded stop time represents the last known accounting activity rather than the time at which the maintenance process executes.

---

# FreeRADIUS Runtime

Each generated configuration tree is intended to be executed by an independent FreeRADIUS instance.

RADIUS Director does not replace the FreeRADIUS distribution.

FreeRADIUS remains responsible for:

- runtime execution
- authentication
- authorization
- accounting
- proxying
- policy execution
- module implementations
- default configuration
- dictionaries
- certificates

RADIUS Director manages only the deployment-specific configuration executed by FreeRADIUS.

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
Managed Configuration Trees
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
- Clear Ownership Boundaries

These principles are further documented in the project's Architectural Decision Records (ADRs).