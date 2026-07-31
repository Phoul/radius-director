# RADIUS Director

RADIUS Director is a declarative deployment platform for FreeRADIUS.

It models multi-tenant RADIUS infrastructure using reusable global objects and tenant-specific configuration, then validates, generates, and deploys managed FreeRADIUS deployments from a single source of truth.

RADIUS Director generates and manages only the deployment-specific configuration required to realize a deployment. The underlying FreeRADIUS distribution, its default configuration, and runtime components remain the responsibility of FreeRADIUS.

## Why RADIUS Director?

Managing multiple FreeRADIUS deployments can be operationally complex. Configuration drift, manual edits, version differences, and deployment inconsistencies make it difficult to provision and maintain RADIUS infrastructure at scale.

RADIUS Director addresses these challenges by treating RADIUS infrastructure as code.

A single declarative configuration defines each tenant's RADIUS deployment. From that configuration, RADIUS Director produces deterministic, version-aware managed configuration that can be deployed consistently across environments.

## Key Features

- Declarative Infrastructure as Code
- Multi-tenant architecture
- One isolated deployment per tenant
- Version-aware FreeRADIUS configuration
- Managed upstream configuration templates
- Deterministic configuration generation
- Managed configuration trees
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
Managed Configuration Tree
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

## Managed Configuration

RADIUS Director manages only the deployment-specific subset of the FreeRADIUS configuration.

Typical managed configuration includes:

- clients.conf
- radiusd.conf
- proxy.conf
- mods-available/
- sites-available/

The managed configuration tree may also include symbolic links required to enable managed modules or virtual servers, such as:

- mods-enabled/
- sites-enabled/

The installed FreeRADIUS distribution remains responsible for:

- dictionaries
- certificates
- policy libraries
- mods-config/
- documentation
- runtime components

This clear ownership boundary minimizes maintenance while remaining closely aligned with upstream FreeRADIUS releases.

## Tenant Model

Each tenant represents an independent FreeRADIUS deployment.

A tenant contains its own:

- Database configuration
- RADIUS server configuration
- NAS assignments
- Managed configuration tree
- Deployment

Each tenant references reusable global objects rather than duplicating shared configuration.

The generated managed configuration for one tenant is completely isolated from every other tenant.

## Version-Aware Deployments

The configured FreeRADIUS version is part of the domain model.

RADIUS Director uses the configured version to:

- Select compatible managed templates
- Validate version-specific configuration
- Render compatible managed configuration
- Select the appropriate deployment implementation

This allows multiple FreeRADIUS versions to be supported while maintaining a consistent configuration model.

## Deployment

Generation and deployment are intentionally separated.

Generation produces a deterministic managed configuration tree.

Deployment is responsible for:

- Installing the managed configuration tree
- Provisioning supporting infrastructure
- Starting or updating FreeRADIUS instances

The deployment layer materializes the managed configuration produced by the generation pipeline. It does not generate or modify configuration itself.

The initial deployment target is Docker-based infrastructure, while allowing future deployment targets without changing the configuration model.

## Project Status

🚧 Early development

Current focus:

- Core object model
- Configuration validation
- Deterministic generation
- Managed configuration rendering
- Deployment architecture
- Initial implementation

## Goals

- Infrastructure as Code
- Declarative configuration
- Deterministic generation
- Multi-tenant deployments
- Version-aware configuration
- Reproducible deployments
- Clear ownership boundaries
- Operational visibility
- Vendor-neutral domain model
- Standards-based configuration

## Non-goals

- Replace FreeRADIUS
- Replace the FreeRADIUS distribution
- Replace existing billing systems
- Become a captive portal
- Become a network management system

## License

Apache License 2.0