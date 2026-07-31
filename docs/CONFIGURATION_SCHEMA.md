# Configuration Schema

This document defines the declarative configuration format used by RADIUS Director.

The schema described here is independent of any implementation language.

---

# Goals

The configuration format should be:

- human readable
- version controllable
- easy to review
- deterministic
- vendor neutral

YAML is currently the preferred configuration format.

---

# Root Configuration

A RADIUS Director configuration consists of a single root configuration document.

The root contains two top-level collections:

- Global Objects
- Tenants

```yaml
global_objects:

  credential_profiles:

  authentication_profiles:

  accounting_profiles:

  monitoring_profiles:

  nas_devices:

tenants:
```

---

# Global Objects

Global Objects are defined once and may be referenced by one or more tenants.

The following Global Object collections are supported:

| Collection | Object |
|------------|--------|
| credential_profiles | Credential Profile |
| authentication_profiles | Authentication Profile |
| accounting_profiles | Accounting Profile |
| monitoring_profiles | Monitoring Profile |
| nas_devices | NAS Device |

Each object is identified by its YAML key.

Example:

```yaml
global_objects:

  credential_profiles:

    default:
      shared_secret: mysecret

  nas_devices:

    mt-core-01.gobcn.ca:
      ip_address: 10.10.10.1
      vendor: mikrotik
```

---

# Tenants

The `tenants` collection contains one or more Tenant objects.

Each tenant is identified by its YAML key.

Example:

```yaml
tenants:

  customer-a:

    database:

    radius_server:

    nas_assignments:
```

Each Tenant owns:

- exactly one Database
- exactly one RADIUS Server
- one or more NAS Assignments

Each tenant represents an independent FreeRADIUS deployment.

Each tenant ultimately produces an independent managed FreeRADIUS configuration tree.

A tenant without any NAS Assignments is considered incomplete and is not valid.

## RADIUS Server

Each Tenant contains one `radius_server` object.

```yaml
radius_server:
  version: 3.2.9
  authentication_port: 1812
  accounting_port: 1813
  coa_port: 3799
```

The `version` field specifies the target FreeRADIUS version for the deployment.

The configured version determines:

- which managed template set is used
- which deployment implementation is selected
- which version-specific validation rules apply

The port properties define the desired listener ports for the deployed FreeRADIUS instance.

The deployment layer is responsible for exposing those ports in the target runtime environment.

---

# Object Collections

Every object collection uses the same pattern.

The YAML key is the object's identifier.

Example:

```yaml
credential_profiles:

  default:
    shared_secret: secret1

  backup:
    shared_secret: secret2
```

The identifiers (`default` and `backup`) become the object identifiers referenced throughout the configuration.

Object identifiers must be unique within their respective collection.

---

# Relationships

Relationship objects reference other objects by identifier.

For example:

```yaml
nas_assignments:

  core-router:

    nas_device: mt-core-01.gobcn.ca

    credential_profile: default

    authentication_profile: default

    accounting_profile: default

    monitoring_profile: default
```

Relationship objects never duplicate configuration owned by other objects.

Referenced objects must exist within the configuration.

---

# Schema Evolution

Configuration compatibility should be versioned.

Breaking changes should include migration guidance.

---

# Validation

Configuration validation should detect:

- invalid YAML structure
- unsupported properties
- missing required properties
- invalid property types
- missing references
- duplicate identifiers
- invalid IP addresses
- invalid object relationships
- unsupported object combinations
- schema version mismatches

Configuration validation should report as many independent errors as practical during a single execution.

No managed configuration should be generated if validation fails.