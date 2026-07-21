# AGENTS.md

## Purpose

This repository contains the source code for **RADIUS Director**, a declarative configuration management system for FreeRADIUS.

This document provides guidance for AI coding assistants contributing to the project.

The project documentation in `docs/` is the authoritative source for architecture and design decisions.

If implementation and documentation disagree, **the documentation takes precedence**.

---

# Project Philosophy

The project values:

- simplicity
- readability
- deterministic behavior
- explicit configuration
- standards over vendor-specific behavior

Avoid unnecessary abstraction or premature optimization.

Implement only what is required today.

---

# Contribution Style

Prefer small, incremental changes.

A single change should ideally:

- implement one feature
- fix one bug
- refactor one component
- improve one area of documentation

Avoid large, sweeping changes affecting multiple subsystems unless explicitly requested.

After completing a task:

- ensure the project builds
- run the test suite
- add or update tests where appropriate
- update documentation if behavior has changed

If a requested feature requires a significant architectural change, stop and explain the trade-offs before implementing it.

---

# Design Principles

## Declarative Configuration

Users describe *what* they want.

The application determines *how* to generate the appropriate FreeRADIUS configuration.

Configuration is the source of truth.

Generated files are disposable artifacts.

---

## Keep the Object Model Small

Do not add new configuration properties without a demonstrated need.

Avoid adding fields simply because they may be useful in the future.

Prefer extending the model later rather than designing speculative features.

---

## Favor Standards

Prefer standard RADIUS behavior whenever practical.

Vendor-specific behavior should be isolated and minimized.

---

## Deterministic Output

The same configuration must always produce identical output.

Never rely on:

- Go map iteration order
- filesystem ordering
- random values
- timestamps

Sort collections before generation whenever ordering matters.

---

## Validation First

Configuration must be validated before generation.

Generation should never attempt to compensate for invalid configuration.

---

# Repository Structure

```
cmd/
    CLI entry point

internal/
    Application implementation

docs/
    Project documentation

examples/
    Example configurations
```

Do not place application logic in `cmd/`.

Keep reusable code inside `internal/`.

---

# Documentation

Before implementing new functionality, consult the relevant documentation.

| Document | Purpose |
|----------|---------|
| VISION.md | Project goals |
| ARCHITECTURE.md | High-level architecture |
| DOMAIN_MODEL.md | Domain concepts |
| OBJECT_MODEL.md | Object definitions |
| CONFIGURATION_SCHEMA.md | Configuration structure |
| VALIDATION.md | Validation rules |
| GENERATOR.md | Generation pipeline |

Do not invent architecture that contradicts these documents.

---

# Coding Standards

Write idiomatic Go.

Prefer:

- small packages
- small functions
- explicit code
- composition over inheritance
- clear naming

Avoid unnecessary interfaces.

Only introduce interfaces when multiple implementations are expected.

Document exported types and functions using standard Go comments.

---

# Error Handling

Return informative errors.

Error messages should help users identify configuration problems.

Include:

- object identifier
- property name
- reason for failure

Prefer messages such as:

```
Tenant "Residential"

NAS Assignment "mt-core-01.gobcn.ca"

Credential Profile "default" does not exist.
```

over generic or ambiguous errors.

When possible, continue validation so multiple configuration errors can be reported in a single run.

---

# Dependencies

Prefer the Go standard library.

Add third-party dependencies only when they provide significant value.

Avoid large frameworks.

Discuss introducing significant new dependencies before implementation.

---

# Testing

New functionality should include tests where practical.

Prefer table-driven tests.

Validation should be thoroughly tested using both valid and invalid example configurations.

Tests should be deterministic and independent.

---

# Generation

Generation should be deterministic.

Generated configuration should be reproducible.

Business logic belongs in validation and the object model—not in templates.

Templates should remain as simple as possible.

Generated files should never require manual modification.

---

# Scope

Do not implement speculative features.

If the documentation does not describe a feature, do not assume it should exist.

When in doubt:

- implement the simplest solution
- leave room for future extension
- avoid unnecessary configuration options

---

# Current Development Priorities

Current implementation order:

1. Project skeleton
2. Configuration parsing
3. Object model
4. Validation
5. CLI (`validate`)
6. Internal generation model
7. Configuration generation
8. Testing

Focus on completing the current stage before introducing functionality from later stages.

---

# Decision Making

When multiple implementation approaches are possible, prefer the solution that is:

1. simpler
2. easier to understand
3. easier to test
4. more deterministic
5. more consistent with the existing documentation

Avoid clever solutions when a straightforward implementation is sufficient.

When uncertain about a design decision, consult the documentation rather than making assumptions.

The goal is to implement the documented design—not redesign the project.