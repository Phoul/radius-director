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

    radius_servers:

    nas_assignments:
```

Each Tenant owns:

- one Database
- one or more RADIUS Servers
- zero or more NAS Assignments

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

---

# Relationships

Relationship Objects reference other objects by identifier.

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

Relationship Objects never duplicate configuration owned by other objects.

---

# Schema Evolution

Configuration compatibility should be versioned.

Breaking changes should include migration guidance.

---

# Validation

Configuration validation should detect:

- missing references
- duplicate identifiers
- invalid IP addresses
- invalid object relationships
- unsupported combinations
- schema version mismatches

No configuration should be generated if validation fails.