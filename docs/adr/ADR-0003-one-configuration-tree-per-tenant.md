# ADR-0003: One Configuration Tree Per Tenant

## Status

Accepted

## Context

RADIUS Director manages multiple independent tenants.

Initially, generation considered producing a single FreeRADIUS configuration containing all tenants.

As the architecture evolved, it became clear that each tenant represents an independent RADIUS deployment with its own infrastructure, including databases, clients, modules, virtual servers, and runtime configuration.

Combining multiple tenants into a single generated configuration would unnecessarily couple otherwise independent deployments.

## Decision

Each tenant SHALL generate its own complete FreeRADIUS configuration tree.

The generated output for each tenant SHALL be entirely independent of every other tenant.

The generation pipeline SHALL iterate over tenants and generate a complete configuration tree for each.

Renderers SHALL operate on a single generated tenant rather than the complete configuration.

## Consequences

Advantages:

- Independent deployment of tenants
- Independent upgrades and maintenance
- Clear separation of tenant infrastructure
- Simplified renderer implementations
- Improved scalability
- Consistent package responsibilities

Trade-offs:

- Some generated files are duplicated across tenants.
- The generator produces multiple configuration trees instead of a single output.

These trade-offs are acceptable because tenant isolation is a primary architectural goal.