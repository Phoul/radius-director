# Glossary

This glossary defines the terminology used throughout the RADIUS Director project.

The purpose of this document is to ensure that documentation, configuration, and implementation use consistent language.

---

# Accounting

The process of recording information about a user's RADIUS session, including session start, stop, duration, and usage.

---

# Accounting Profile

A reusable object that defines how accounting is handled for one or more NAS devices.

Accounting Profiles are referenced by NAS Assignments.

Examples include:

- accounting storage
- interim update handling
- session cleanup
- retention policies

---

# Authentication

The process of verifying a user's identity and determining whether access should be granted.

---

# Authentication Profile

A reusable object that defines authentication behaviour.

Authentication Profiles are referenced by NAS Assignments.

Examples include:

- simultaneous use
- SQL policies
- session verification
- vendor-specific behaviour

---

# CoA (Change of Authorization)

A RADIUS message that modifies an active session without disconnecting the user.

CoA is defined by RFC 5176.

---

# Credential Profile

A reusable object containing shared credentials referenced by one or more NAS Assignments.

Typical credentials include:

- RADIUS shared secret
- CoA shared secret

---

# Domain Model

The collection of objects used to describe a RADIUS deployment.

The domain model is the primary source of truth for configuration generation.

---

# FreeRADIUS

The open-source RADIUS server that executes authentication, authorization, accounting, proxying, and policy logic.

RADIUS Director generates configuration for FreeRADIUS but does not replace it.

---

# Generator

The component responsible for transforming the domain model into FreeRADIUS configuration files.

Generated configuration should be deterministic and reproducible.

---

# Global Object

A reusable object defined once within a RADIUS Director installation.

Global Objects may be referenced by one or more tenants and should not contain tenant-specific configuration.

Examples include:

- Credential Profiles
- Authentication Profiles
- Accounting Profiles
- Monitoring Profiles
- NAS Devices

---

# Infrastructure as Code (IaC)

The practice of managing infrastructure through version-controlled configuration rather than manual changes.

RADIUS Director adopts Infrastructure as Code principles wherever practical.

---

# Monitoring Profile

A reusable object that defines operational monitoring for one or more NAS devices.

Monitoring Profiles are referenced by NAS Assignments.

Examples include:

- connectivity tests
- authentication tests
- CoA tests
- SNMP verification

---

# NAS Device

A reusable Global Object representing a physical or virtual RADIUS client.

A NAS Device describes the device itself rather than how it is used.

Typical properties include:

- Name
- Address
- Vendor
- Model
- Description

Operational behaviour such as credentials, authentication, accounting, and monitoring is defined through NAS Assignments.

---

# NAS Assignment

A Relationship Object that describes how a tenant uses a NAS Device.

A NAS Assignment associates a NAS Device with one or more operational profiles, including:

- Credential Profile
- Authentication Profile
- Accounting Profile
- Monitoring Profile

Multiple tenants may reference the same NAS Device while maintaining independent NAS Assignments.

---

# Proxying

The process of forwarding RADIUS requests to another RADIUS server.

Proxying may be used for roaming, federation, or vendor-specific architectures.

---

# RADIUS Director

The software project described by this repository.

RADIUS Director manages the configuration, validation, deployment, and operation of standards-based RADIUS infrastructure.

It is not a RADIUS server.

---

# RADIUS Server

A running instance of FreeRADIUS that consumes generated configuration.

A deployment may contain one or more RADIUS servers.

---

# Relationship Object

An object that describes how Global Objects are composed into a tenant deployment.

Relationship Objects belong to a tenant.

---

# Source of Truth

The authoritative representation of the desired system state.

In RADIUS Director, the source of truth is the declarative configuration model stored in version control.

Generated configuration is not the source of truth.

---

# Standards-Based

A design philosophy that models open standards before vendor-specific implementations.

Examples include RFC 2865, RFC 2866, and RFC 5176.

---

# Tenant

An independent RADIUS deployment managed by RADIUS Director.

Each tenant owns its tenant-specific infrastructure, including:

- Database
- RADIUS Servers
- NAS Assignments

Tenants reference reusable Global Objects rather than duplicating them.

---

# Tenant Object

An object owned exclusively by a tenant.

Tenant Objects define tenant-specific infrastructure.

---

# Validation

The process of verifying that a configuration is internally consistent before generation or deployment.

Validation should detect errors such as:

- missing references
- duplicate identifiers
- invalid IP addresses
- schema violations
- unsupported relationships

Configuration generation should not proceed if validation fails.

---

# Vendor Neutral

A design principle that avoids modelling specific vendor products unless required for standards compliance or device capabilities.

The goal is to maximize portability across network vendors.

---

# YAML

The human-readable configuration format currently planned for describing the RADIUS Director domain model.

YAML files represent the source of truth for generated configuration.