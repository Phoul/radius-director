# Validation

This document defines the validation rules applied to the RADIUS Director configuration model.

Validation ensures that configuration is internally consistent before any FreeRADIUS configuration is generated.

Generation must not proceed if validation fails.

---

# Validation Goals

Validation should:

- detect configuration errors as early as possible
- produce clear and actionable error messages
- prevent generation of invalid FreeRADIUS configuration
- validate the configuration model rather than generated files

---

# Validation Stages

Validation is performed in several stages.

## 1. Schema Validation

Verifies that the configuration is syntactically valid.

Examples include:

- required properties
- supported property names
- property types
- valid YAML structure

---

## 2. Object Validation

Each object validates its own properties.

Examples include:

- required fields
- valid IP addresses
- supported database engines
- supported vendors

---

## 3. Reference Validation

Ensures that all object references exist.

Examples include:

- NAS Assignment references an existing NAS Device
- NAS Assignment references an existing Credential Profile
- NAS Assignment references an existing Authentication Profile
- NAS Assignment references an existing Accounting Profile
- NAS Assignment references an existing Monitoring Profile

---

## 4. Relationship Validation

Ensures that object relationships are valid.

Examples include:

- duplicate object identifiers
- duplicate IP addresses
- duplicate NAS Assignments
- unsupported object combinations

---

## 5. Generation Validation

Verifies that sufficient information exists to generate configuration.

Examples include:

- every tenant has a database
- every tenant has at least one RADIUS Server
- every NAS Assignment is complete

---

# Object Validation Rules

## Credential Profile

Validation rules:

- identifier must be unique
- shared_secret must be specified

---

## Authentication Profile

Validation rules:

- identifier must be unique

Additional validation is implementation specific.

---

## Accounting Profile

Validation rules:

- identifier must be unique

Additional validation is implementation specific.

---

## Monitoring Profile

Validation rules:

- identifier must be unique

Additional validation is implementation specific.

---

## NAS Device

Validation rules:

- identifier must be unique
- ip_address must be a valid IPv4 or IPv6 address
- ip_address must be unique
- vendor must be specified

---

## Tenant

Validation rules:

- identifier must be unique
- exactly one Database must be defined
- at least one RADIUS Server must be defined
- at least one NAS Assignment must be defined

---

## Database

Validation rules:

- engine must be supported
- host must be specified
- port must be valid
- database must be specified
- username must be specified
- password must be specified

---

## RADIUS Server

Validation rules:

- identifier must be unique

Additional validation is implementation specific.

---

## NAS Assignment

Validation rules:

- identifier must be unique
- referenced NAS Device must exist
- referenced Credential Profile must exist
- referenced Authentication Profile must exist
- referenced Accounting Profile must exist
- referenced Monitoring Profile must exist

---

# Validation Errors

Validation errors should:

- identify the object containing the error
- identify the property causing the error
- describe the problem clearly
- provide enough information to locate the error

Example:

```
ERROR

Tenant: residential

NAS Assignment: mt-core-01.gobcn.ca

Credential Profile "default" does not exist.
```

Validation should continue where practical so that multiple errors can be reported in a single run.

---

# Validation Philosophy

Validation should be strict.

Configuration that cannot be generated correctly should be rejected.

The validator should fail early, fail clearly, and never generate configuration that is known to be invalid.