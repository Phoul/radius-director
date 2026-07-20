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

# High-Level Structure

```yaml
platform:

tenants:

credential_profiles:

authentication_profiles:

accounting_profiles:

monitoring_profiles:

nas:

databases:

radius_servers:
```

Detailed schema definitions will be added as the domain model is finalized.

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