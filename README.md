# RADIUS Director

RADIUS Director is a declarative deployment platform for FreeRADIUS.

It models multi-tenant RADIUS infrastructure using reusable global objects and tenant-specific configuration, then validates, generates, and deploys complete FreeRADIUS instances from a single source of truth.

## Why RADIUS Director?

Managing multiple FreeRADIUS deployments can be operationally complex. Configuration drift, manual edits, version differences, and deployment inconsistencies make it difficult to provision and maintain RADIUS infrastructure at scale.

RADIUS Director addresses these challenges by treating RADIUS infrastructure as code.

A single declarative configuration defines each tenant's RADIUS environment. From that configuration, RADIUS Director produces deterministic, version-aware FreeRADIUS deployments.

## Key Features

- Declarative Infrastructure as Code
- Multi-tenant architecture
- One isolated FreeRADIUS deployment per tenant
- Version-aware FreeRADIUS configuration
- Managed upstream configuration templates
- Deterministic configuration generation
- Integrated deployment pipeline
- Configuration validation before deployment
- Reproducible infrastructure
- Vendor-neutral configuration model

## Architecture

RADIUS Director separates responsibilities into distinct stages.

```text
Configuration
        │
        ▼
Validation
        │
        ▼
Generator
        │
        ▼
Renderer
        │
        ▼
Output
        │
        ▼
Writer
        │
        ▼
Deployment
        │
        ▼
Running FreeRADIUS
```

Each stage has a single responsibility and performs no work belonging to another stage.

## Tenant Model

Each tenant represents an independent FreeRADIUS deployment.

A tenant contains its own:

- RADIUS server configuration
- SQL configuration
- Clients
- Policies
- Generated configuration tree
- Deployment

The generated configuration for one tenant is completely isolated from every other tenant.

## Version-Aware Deployments

The configured FreeRADIUS version is part of the domain model.

RADIUS Director uses the configured version to:

- Select compatible configuration templates
- Validate version-specific configuration
- Generate compatible FreeRADIUS configuration
- Select the appropriate deployment image

This allows multiple FreeRADIUS versions to be supported while maintaining a consistent configuration model.

## Deployment

RADIUS Director is responsible for deploying complete FreeRADIUS environments, not merely generating configuration files.

Deployment responsibilities include:

- Deploying FreeRADIUS
- Selecting version-specific templates
- Creating runtime configuration
- Managing supporting infrastructure
- Producing reproducible deployments

The initial deployment target is Docker-based infrastructure, with future deployment targets remaining possible without changing the configuration model.

## Project Status

🚧 Early design phase

Current focus:

- Domain model
- Configuration schema
- Generator architecture
- Deployment architecture
- Documentation
- Initial implementation

## Goals

- Infrastructure as Code
- Declarative configuration
- Deterministic generation
- Multi-tenant deployments
- Version-aware configuration
- Reproducible deployments
- Operational visibility
- Vendor-neutral domain model
- Standards-based configuration

## Non-goals

- Replace FreeRADIUS
- Replace existing billing systems
- Become a captive portal
- Become a network management system

## License

Apache License 2.0