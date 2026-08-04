# Object Model

This document defines each object in the RADIUS Director domain model.

Every object describes a single operational concept and follows the same structure:

- Purpose
- Ownership
- Properties
- Relationships
- Validation
- Generation

The object model serves as the blueprint for configuration validation and managed FreeRADIUS configuration generation.

---

# General Principles

## Object Identifiers

Every object is identified by its key within the configuration.

For example:

```yaml
nas_devices:

  mt-core-01.gobcn.ca:
    ip_address: 10.10.10.1
    vendor: mikrotik
```

The key (`mt-core-01.gobcn.ca`) is the object's identifier.

Object identifiers:

- must be unique within their collection
- are referenced by other objects
- should remain stable over time

---

# Credential Profile

## Purpose

Defines the shared RADIUS credentials used when communicating with a NAS Device.

## Ownership

Global Object

## Properties

- shared_secret

## Relationships

Referenced by one or more NAS Assignments.

## Validation

- shared_secret must be present

## Generation

Used when rendering managed client definitions and CoA configuration.

---

# Authentication Profile

## Purpose

Defines authentication behaviour.

## Ownership

Global Object

## Properties

Initially implementation-specific.

Examples may include:

- simultaneous use
- SQL policies
- session verification

## Relationships

Referenced by one or more NAS Assignments.

## Validation

Profile-specific validation rules.

## Generation

Used when generating authentication policies.

---

# Accounting Profile

## Purpose

Defines accounting behaviour.

## Ownership

Global Object

## Properties

Initially implementation-specific.

Examples may include:

- accounting storage
- interim update handling
- session cleanup

## Relationships

Referenced by one or more NAS Assignments.

## Validation

Profile-specific validation rules.

## Generation

Used when generating accounting configuration.

---

# Monitoring Profile

## Purpose

Defines operational monitoring.

## Ownership

Global Object

## Properties

Initially implementation-specific.

Examples may include:

- authentication tests
- CoA tests
- SNMP verification

## Relationships

Referenced by one or more NAS Assignments.

## Validation

Profile-specific validation rules.

## Generation

Used when generating monitoring configuration.

---

# NAS Device

## Purpose

Represents a physical or virtual RADIUS client.

## Ownership

Global Object

## Properties

- ip_address
- vendor

## Relationships

Referenced by one or more NAS Assignments.

## Validation

- ip_address must be a valid IP address
- vendor must be supported

## Generation

Used when rendering managed FreeRADIUS client definitions.

---

# Trusted RADIUS Client

## Purpose

Represents a trusted RADIUS client that communicates with FreeRADIUS but is not itself a NAS Device.

Examples include:

- billing systems
- provisioning systems
- management platforms
- monitoring systems
- other RADIUS servers

Trusted RADIUS Clients participate only in client authentication and authorization.

They are not used when generating CoA proxy configuration.

## Ownership

Global Object

## Properties

- ip_address

## Relationships

Referenced by one or more Trusted RADIUS Client Assignments.

## Validation

- ip_address must be a valid IP address

## Generation

Used when rendering managed FreeRADIUS client definitions.

Trusted RADIUS Clients do not participate in managed CoA configuration generation.

---

# Tenant

## Purpose

Represents an independent RADIUS deployment.

## Ownership

Top-level object.

## Properties

None.

A tenant is composed of other Tenant Objects.

## Relationships

Contains:

- Database
- RADIUS Server
- NAS Assignments
- Trusted RADIUS Client Assignments

Each tenant represents a complete, independent FreeRADIUS deployment.

The tenant's RADIUS Server object defines the target FreeRADIUS version for the deployment.

## Validation

- must contain exactly one Database
- must contain exactly one RADIUS Server

## Generation

Defines the scope of the managed configuration generated for the tenant.

Each tenant produces an independent managed FreeRADIUS configuration tree.

---

# Database

## Purpose

Defines the primary database used by a tenant.

## Ownership

Tenant Object

## Properties

- engine
- host
- port
- database
- username
- password

## Relationships

Owned by a Tenant.

## Validation

- host must be specified
- database must be specified
- credentials must be valid

## Generation

The deployment layer is responsible for determining how the database is connected to at runtime.

---

# RADIUS Server

## Purpose

Represents the desired runtime characteristics of a FreeRADIUS deployment.

## Ownership

Tenant Object

## Properties

- version
- authentication_port
- accounting_port
- coa_port

## Relationships

Owned by a Tenant.

## Validation

- version must be specified
- version must be supported by the current RADIUS Director release
- authentication_port must be between 1 and 65535
- accounting_port must be between 1 and 65535
- coa_port must be between 1 and 65535

## Generation

The RADIUS Server defines the runtime characteristics of the generated deployment.

The configured version determines:

- which managed template set is used
- which deployment image is selected
- which version-specific validation rules apply

The port properties define the desired listener ports for the deployed FreeRADIUS instance.

The deployment layer is responsible for exposing those ports in the target runtime environment.

---

# NAS Assignment

## Purpose

Defines how a tenant uses a NAS Device.

## Ownership

Relationship Object

## Properties

- nas_device
- credential_profile
- authentication_profile
- accounting_profile
- monitoring_profile

## Relationships

References:

- NAS Device
- Credential Profile
- Authentication Profile
- Accounting Profile
- Monitoring Profile

Owned by a Tenant.

## Validation

- nas_device must be specified
- credential_profile must be specified
- authentication_profile must be specified
- accounting_profile must be specified
- monitoring_profile must be specified
- every referenced object must exist
- duplicate assignments are not permitted

## Generation

Combines Global Objects into tenant-specific managed configuration used to render the tenant's managed FreeRADIUS configuration tree.

---

# Trusted RADIUS Client Assignment

## Purpose

Defines how a tenant uses a Trusted RADIUS Client.

Unlike NAS Assignments, Trusted RADIUS Client Assignments describe only the credentials required for a trusted client to communicate with the tenant's FreeRADIUS deployment.

## Ownership

Relationship Object

## Properties

- trusted_radius_client
- credential_profile

## Relationships

References:

- Trusted RADIUS Client
- Credential Profile

Owned by a Tenant.

## Validation

- trusted_radius_client must be specified
- credential_profile must be specified
- every referenced object must exist
- duplicate assignments are not permitted

## Generation

Combines Trusted RADIUS Clients with Credential Profiles to render managed client definitions.

Trusted RADIUS Client Assignments do not participate in CoA proxy configuration generation.