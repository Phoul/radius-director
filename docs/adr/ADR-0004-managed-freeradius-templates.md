# ADR-0004: Managed FreeRADIUS Templates

## Status

Accepted

## Context

Many FreeRADIUS configuration files are large, complex, and evolve between FreeRADIUS releases.

Examples include:

- mods-available/eap
- mods-available/ldap
- sites-enabled/default

Generating every line of every configuration file in Go would require RADIUS Director to duplicate large portions of the upstream FreeRADIUS configuration and continuously track upstream changes.

## Decision

RADIUS Director SHALL generate configuration from managed FreeRADIUS templates.

Templates SHALL be based on the upstream FreeRADIUS configuration for a specific supported FreeRADIUS release.

Templates SHALL be versioned alongside the supported FreeRADIUS version.

The configured FreeRADIUS version determines which template set is used during generation.

Templates SHALL contain placeholders only for values managed by the domain model.

Renderers are responsible for supplying tenant-specific values.

Static configuration that is not managed by RADIUS Director SHALL remain in the template.

## Consequences

Advantages:

- Remains closely aligned with upstream FreeRADIUS
- Supports multiple FreeRADIUS versions simultaneously
- Templates evolve independently for each supported release
- Easier upgrades to newer FreeRADIUS versions
- Smaller renderer implementations
- Reduced maintenance burden
- Only managed values are generated

Trade-offs:

- Template files become part of the project.
- Upstream template changes must occasionally be incorporated.

These trade-offs are preferable to reimplementing large FreeRADIUS configuration files in code.