# RADIUS Director

RADIUS Director is a declarative deployment platform for FreeRADIUS.

It models multi-tenant RADIUS infrastructure using reusable global objects and tenant-specific configuration, then validates and generates managed FreeRADIUS deployments from a single source of truth.

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
- Per-tenant deployment manifests
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

Generation and deployment are intentionally separate. RADIUS Director generates the configuration and deployment metadata; the deployment layer is responsible for materializing that configuration into a running FreeRADIUS environment.

## Managed Configuration

RADIUS Director manages only the deployment-specific subset of the FreeRADIUS configuration.

Typical managed configuration includes:

- clients.conf
- radiusd.conf
- proxy.conf
- clients.d/
- proxy.d/
- mods-available/
- mods-enabled/
- sites-available/
- sites-enabled/
- mods-config/files/authorize

The managed configuration tree may also include symbolic links required to enable managed modules or virtual servers, such as:

- mods-enabled/
- sites-enabled/

The installed FreeRADIUS distribution remains responsible for configuration and runtime components that are not managed by RADIUS Director, including:

- dictionaries
- certificates
- policy libraries
- base/default configuration supplied by the FreeRADIUS distribution
- runtime components
- other upstream files not explicitly managed by a tenant

This ownership boundary allows RADIUS Director to remain closely aligned with upstream FreeRADIUS releases without requiring the entire FreeRADIUS distribution to be maintained as part of the project.

## Tenant Model

Each tenant represents an independent FreeRADIUS deployment.

A tenant contains its own:

- Database configuration
- RADIUS server configuration
- NAS assignments
- Managed configuration tree
- Deployment profile
- Deployment manifest

Each tenant references reusable global objects rather than duplicating shared configuration.

The generated managed configuration for one tenant is completely isolated from every other tenant.

## Configuration

The primary configuration is a YAML document containing global objects and tenants.

An example configuration is provided in:

```text
examples/example.yaml
```

A simplified example looks like:

```yaml
global_objects:
  credential_profiles:
    default:
      shared_secret: shared-secret

  authentication_profiles:
    default:
      simultaneous_use: 1

  accounting_profiles:
    default: {}

  monitoring_profiles:
    default: {}

  deployment_profiles:
    default:
      template: default
      overlays: []

  nas_devices:
    mt-core-01.gobcn.ca:
      ip_address: 10.10.10.1
      vendor: mikrotik

tenants:
  customer-a:
    authentication_profile: default
    deployment_profile: default

    database:
      engine: mysql
      host: db.example.com
      port: 3306
      database: radius
      username: radius
      password: password

    radius_server:
      version: 3.2.10
      authentication_port: 1812
      accounting_port: 1813
      coa_port: 3799

    nas_assignments:
      core-router:
        nas_device: mt-core-01.gobcn.ca
        credential_profile: default
        accounting_profile: default
        monitoring_profile: default
```

The complete example in `examples/example.yaml` should be used as the reference for the currently supported configuration structure.

## Version-Aware Deployments

The configured FreeRADIUS version is part of the domain model.

For example:

```yaml
radius_server:
  version: 3.2.10
```

RADIUS Director uses the configured version to:

- Select compatible managed templates
- Validate version-specific configuration
- Render compatible managed configuration
- Select the appropriate deployment implementation

This allows multiple FreeRADIUS versions to be supported while maintaining a consistent configuration model.

The eventual Docker deployment uses the configured FreeRADIUS version to select the corresponding official `freeradius/freeradius-server` image.

## Deployment Profiles

Deployment profiles define how a tenant's FreeRADIUS configuration is assembled.

A deployment profile currently specifies a template and optional overlays:

```yaml
global_objects:
  deployment_profiles:
    default:
      template: default
      overlays: []
```

Deployment profiles may also specify paths that should be removed from the effective FreeRADIUS configuration when the deployment is materialized:

```yaml
global_objects:
  deployment_profiles:
    default:
      template: default
      overlays: []
      remove:
        - sites-enabled/inner-tunnel
```

The removal paths are recorded in the tenant's deployment manifest:

```text
.radius-director/manifest.yaml
```

For example:

```yaml
remove:
    - sites-enabled/inner-tunnel
```

The generated tenant directory itself is not modified by these removal instructions.

Instead, the deployment layer is responsible for applying the removals to the effective FreeRADIUS configuration inside the deployment environment. This allows a deployment to remove files supplied by the underlying FreeRADIUS distribution without deleting them from the generated tenant configuration on the host.

## Usage

RADIUS Director is currently run from the repository root using Go's `go run` command.

### Validate a configuration

To validate a configuration:

```powershell
go run ./cmd/radius-director validate ./examples/example.yaml
```

A successful validation indicates that the configuration is valid for the current RADIUS Director configuration model.

### Generate tenant configurations

To generate the managed configuration for all tenants:

```powershell
go run ./cmd/radius-director generate ./examples/example.yaml ./generated
```

The output directory will contain a separate directory for each tenant.

For example:

```text
generated/
└── customer-a/
    ├── .radius-director/
    │   └── manifest.yaml
    ├── clients.conf
    ├── clients.d/
    │   └── radius-director.conf
    ├── mods-available/
    │   └── sql
    ├── mods-config/
    │   └── files/
    │       └── authorize
    ├── mods-enabled/
    │   └── sql
    ├── proxy.conf
    ├── proxy.d/
    │   └── radius-director.conf
    ├── sites-available/
    │   ├── coa
    │   └── default
    ├── sites-enabled/
    │   ├── coa -> ../sites-available/coa
    │   └── default -> ../sites-available/default
    └── users
```

The generated tenant directory is the managed configuration artifact for that tenant.

It should generally be treated as generated output rather than manually edited configuration.

### Command reference

```text
radius-director validate <config.yaml>
radius-director generate <config.yaml> <output-directory>
radius-director maintenance accounting <config.yaml> <tenant>
```

For command-line help:

```powershell
go run ./cmd/radius-director --help
```

## Development

RADIUS Director is written in Go.

Run the complete test suite from the repository root:

```powershell
go test ./...
```

Run the CLI without installing it:

```powershell
go run ./cmd/radius-director --help
```

When testing generation, write output to a separate directory rather than modifying the repository's template directories.

For example:

```powershell
go run ./cmd/radius-director generate ./examples/example.yaml ./generated-test
```

The generated directory can then be inspected before being removed:

```powershell
Remove-Item -Recurse -Force ./generated-test
```

## Deployment

Generation and deployment are intentionally separated.

Generation produces a deterministic managed configuration tree and a per-tenant deployment manifest.

Deployment is responsible for:

- Installing or materializing the managed configuration
- Provisioning supporting infrastructure
- Applying deployment-manifest instructions
- Starting or updating FreeRADIUS instances

The initial deployment target is Docker-based infrastructure, while allowing future deployment targets without changing the configuration model.

### Docker deployment architecture

The Docker deployment is currently under development.

The intended architecture uses the official FreeRADIUS Docker image rather than maintaining a custom FreeRADIUS base image.

Each tenant will have its own FreeRADIUS container, using the FreeRADIUS version specified by that tenant's configuration.

Conceptually:

```text
RADIUS Director YAML
        │
        ▼
   Configuration
        │
        ▼
   Per-tenant generated configuration
        │
        ├── customer-a/
        └── customer-b/
        │
        ▼
   Docker deployment
        │
        ├── official FreeRADIUS image
        │       + tenant configuration
        │       + tenant deployment manifest
        │
        └── one isolated container per tenant
```

The deployment entrypoint is responsible for constructing the effective `/etc/freeradius` configuration inside the container.

The generated tenant configuration is mounted separately from the container's effective FreeRADIUS configuration. Deployment-manifest removal instructions are applied inside the container so that files supplied by the base FreeRADIUS image can be removed without modifying the generated configuration on the Docker host.

## Project Status

🚧 Early development

Currently implemented:

- Core object model
- Configuration loading and validation
- Deterministic generation
- Managed configuration rendering
- Per-tenant managed configuration trees
- Per-tenant deployment manifests
- Deployment-profile removal instructions
- Initial Docker deployment entrypoint
- Automated test coverage for core generation and writing functionality

Current development focus:

- Docker deployment architecture
- Docker Compose generation
- Per-tenant container deployment
- Integration with official FreeRADIUS Docker images
- Deployment and runtime lifecycle

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