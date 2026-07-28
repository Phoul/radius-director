# AGENTS.md

## Purpose

This repository contains the source code for **RADIUS Director**, a declarative deployment platform for FreeRADIUS.

RADIUS Director models multi-tenant RADIUS infrastructure using a declarative domain model, validates that model, generates version-aware FreeRADIUS configuration, and deploys complete FreeRADIUS environments from a single source of truth.

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

# Current Architecture

The implementation pipeline is:

```
Configuration
    ↓
Validation
    ↓
Generator
    ↓
Renderer
    ↓
Output
    ↓
Writer
    ↓
Deployment
```

Each stage has a single responsibility.

## Validation

Validation is responsible for:

- object validation
- reference validation
- relationship validation

Validation should report as many configuration errors as practical.

Generation assumes validation has already succeeded.

Validation must never modify configuration.

---

## Generator

The generator performs a pure transformation from the validated configuration model into an internal generation model.

The generator:

- performs no validation
- performs no filesystem operations
- performs no rendering
- contains no side effects

The generated model should already be deterministic.

Collections should be sorted here whenever ordering matters.

---

## Renderer

Renderers convert the generation model into file contents.

Renderers:

- perform no validation
- perform no filesystem operations
- perform no business logic
- do not reference the original configuration model

Prefer straightforward Go code over text templates unless templates provide a clear advantage.

---

## Output

The output package assembles rendered files into an `Output` object.

It does not perform rendering itself.

---

## Writer

The writer performs filesystem operations only.

It:

- creates directories
- writes files
- overwrites existing files
- rejects path traversal

It must not:

- delete files
- validate configuration
- render templates
- contain FreeRADIUS-specific logic

---

## Deployment

The deployment package is responsible for transforming generated configuration into a runnable FreeRADIUS environment.

Deployment responsibilities include:

- selecting the target FreeRADIUS version
- selecting version-compatible managed templates
- generating deployment artifacts
- deploying supporting infrastructure
- deploying and updating FreeRADIUS instances

Deployment must not:

- perform validation
- modify generated configuration
- render configuration
- contain business logic

The deployment package should remain independent of the generator and renderer.

Docker is the initial deployment target, but the deployment architecture should remain sufficiently abstract to support additional deployment targets in the future.

---

## Package Responsibilities

Packages should remain independent.

For example:

- Validation must not render.
- Generator must not validate.
- Renderer must not generate.
- Output must not perform rendering.
- Writer must not know anything about FreeRADIUS.
- Deployment must not perform generation or rendering.

Keep responsibilities narrow and explicit.

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

Prefer table-driven tests where they improve readability.

Tests should be deterministic and independent.

Validation should be thoroughly tested using both valid and invalid example configurations.

When adding a new package, include comprehensive unit tests before moving to the next stage.

Prefer testing observable behaviour rather than implementation details.

---

# Generation

Generation should be deterministic.

Generated configuration should be reproducible.

Business logic belongs in validation and the generator.

Renderers should contain only formatting logic.

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
8. Writer
9. Deployment
10. Testing

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

---

# Important

When implementing new functionality:

- preserve the existing architecture
- avoid moving responsibilities between packages
- avoid introducing interfaces unless there is a demonstrated need
- avoid speculative abstractions
- if a requested change appears to conflict with this document or the project documentation, explain the conflict instead of making architectural changes