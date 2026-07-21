# Generator

This document describes how RADIUS Director transforms the declarative configuration model into a complete FreeRADIUS configuration.

The generator consumes a validated configuration model and produces deterministic FreeRADIUS configuration files.

---

# Design Goals

The generator should be:

- deterministic
- repeatable
- idempotent
- vendor neutral
- standards based

Given the same input configuration, the generator must always produce identical output.

---

# Generation Pipeline

Configuration generation consists of several stages.

```
Read Configuration
        │
        ▼
Parse YAML
        │
        ▼
Build Object Model
        │
        ▼
Resolve References
        │
        ▼
Validate
        │
        ▼
Generate Internal Model
        │
        ▼
Generate FreeRADIUS Configuration
        │
        ▼
Write Output Files
```

Each stage must complete successfully before the next stage begins.

---

# Stage 1 - Read Configuration

The generator reads the declarative configuration from one or more YAML files.

No validation is performed at this stage beyond confirming that the files can be read.

---

# Stage 2 - Parse Configuration

The YAML configuration is parsed into the internal object model.

This stage converts configuration into strongly typed objects.

No FreeRADIUS configuration is generated during parsing.

---

# Stage 3 - Resolve References

Object references are resolved.

Examples include:

- NAS Assignments locating NAS Devices
- NAS Assignments locating Credential Profiles
- Tenants locating Databases
- Tenants locating RADIUS Servers

After this stage, the object graph is fully connected.

---

# Stage 4 - Validate

The complete object model is validated.

Generation must stop immediately if validation fails.

Validation rules are described in VALIDATION.md.

---

# Stage 5 - Build Internal Model

The validated object graph is transformed into an internal representation optimized for configuration generation.

This stage may combine multiple objects into structures that more closely represent the generated FreeRADIUS configuration.

The internal model is an implementation detail and is not part of the public configuration schema.

---

# Stage 6 - Generate Configuration

The internal model is used to generate FreeRADIUS configuration.

Generation should be independent for each configuration component wherever practical.

Typical generated files may include:

- clients.conf
- proxy.conf
- mods-enabled/sql
- policy.d
- sites-enabled

Additional generated files may be added in future versions.

---

# Stage 7 - Write Output

Generated files are written to the configured output directory.

Existing generated files may be replaced.

Files not managed by the generator should not be modified.

---

# Deterministic Generation

Configuration generation must be deterministic.

The generator should not depend on:

- object insertion order
- filesystem ordering
- map iteration order
- random values

Objects should be generated in a predictable order whenever practical.

---

# Error Handling

Generation must stop if:

- configuration cannot be parsed
- references cannot be resolved
- validation fails
- generation encounters an unrecoverable error

No partial configuration should be written.

---

# Source of Truth

The declarative configuration is the source of truth.

Generated FreeRADIUS configuration is an artifact that may be safely regenerated at any time.

Generated configuration should never be edited manually.

---

# Design Philosophy

The generator should remain as simple as possible.

Business logic belongs in the configuration model.

Validation belongs in the validator.

The generator's responsibility is to transform a valid object model into deterministic FreeRADIUS configuration.