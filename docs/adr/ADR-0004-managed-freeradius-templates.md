# ADR-0004: Managed FreeRADIUS Templates

## Status

Accepted

## Context

Many FreeRADIUS configuration files are large, complex, and evolve between FreeRADIUS releases.

Examples include:

- mods-available/eap
- mods-available/ldap
- sites-available/default

Generating every line of every configuration file in Go would require RADIUS Director to duplicate large portions of the upstream FreeRADIUS configuration and continuously track upstream changes.

FreeRADIUS configuration syntax, directives, defaults, and implementation-specific values may also differ between releases.

The RADIUS Director domain model should remain stable and describe operational concepts rather than expose FreeRADIUS-version-specific implementation details.

For example, the domain model may describe a NAS Device as:

```yaml
vendor: mikrotik
```

while a particular FreeRADIUS release may require a corresponding configuration value such as:

```text
nas_type = mikrotik_snmp
```

Such FreeRADIUS-specific mappings should not need to become properties of the domain model when they can be deterministically derived from existing domain information.

## Decision

RADIUS Director SHALL generate managed configuration from version-specific FreeRADIUS templates.

Managed templates SHALL be based on the upstream FreeRADIUS configuration for a specific supported FreeRADIUS release.

Templates SHALL be versioned alongside the supported FreeRADIUS version.

The configured FreeRADIUS version determines which template set is used during generation.

Templates SHALL act as the version-specific translation layer between the RADIUS Director domain model and FreeRADIUS configuration.

Renderers SHALL supply domain-model data to templates.

Templates MAY translate domain-model values into FreeRADIUS-specific configuration values where the mapping is deterministic.

Version-specific FreeRADIUS implementation details SHOULD remain within the corresponding template set rather than being introduced into the domain model or hard-coded into renderer logic.

For example, a template may translate:

```yaml
vendor: mikrotik
```

into:

```text
nas_type = mikrotik_snmp
```

when that is the appropriate representation for the FreeRADIUS version associated with that template set.

If a later FreeRADIUS version changes or removes that representation, only the corresponding version-specific template should need to change.

Templates SHALL contain placeholders only for values represented by the domain model.

New domain-model properties SHALL NOT be introduced solely to expose FreeRADIUS implementation details when those details can be deterministically derived by the version-specific template.

Only configuration files that are managed by RADIUS Director are represented as templates.

Within those managed templates, static configuration that is not controlled by or derived from the domain model SHALL remain unchanged from the corresponding upstream FreeRADIUS release.

Configuration that is outside the managed configuration boundary remains the responsibility of the installed FreeRADIUS distribution and is not templated by RADIUS Director.

## Responsibilities

The domain model is responsible for describing deployment intent using stable operational concepts.

The renderer is responsible for:

- selecting the template set associated with the configured FreeRADIUS version
- discovering the managed templates contained within the selected template set
- supplying the appropriate domain-model data to each template
- executing each template
- preserving the template directory structure within the generated managed configuration tree
- returning the rendered managed configuration

The template is responsible for:

- representing the managed FreeRADIUS configuration for its specific FreeRADIUS release
- substituting values supplied by the domain model
- translating domain-model concepts into version-specific FreeRADIUS configuration where required
- preserving upstream configuration that is not managed or derived by RADIUS Director

Version-specific FreeRADIUS configuration knowledge SHOULD therefore reside in templates wherever practical rather than in renderer or domain-model code.

## Consequences

Advantages:

- Keeps the domain model independent of FreeRADIUS-version-specific implementation details
- Keeps version-specific configuration knowledge with the corresponding FreeRADIUS release
- Remains closely aligned with upstream FreeRADIUS
- Supports multiple FreeRADIUS versions simultaneously
- Templates evolve independently for each supported release
- Allows FreeRADIUS syntax and implementation details to change without requiring domain-model changes
- Reduces version-specific logic in Go
- Easier upgrades to newer FreeRADIUS versions
- Smaller renderer implementations
- Reduced maintenance burden
- Only managed or derived values differ from upstream configuration
- Adding managed configuration files typically requires only adding a new template to the appropriate version-specific template set.

Trade-offs:

- Template files become part of the project
- Templates may contain limited conditional or translation logic
- Upstream template changes must occasionally be incorporated
- Version-specific mappings must be reviewed when adding support for a new FreeRADIUS release

Template logic SHOULD remain simple.

Complex business logic belongs in the domain model or generator rather than in templates. If version-specific template logic becomes complex or cannot be expressed as a deterministic translation of the domain model, the architecture should be reconsidered rather than allowing templates to accumulate application logic.

These trade-offs are preferable to reimplementing large portions of FreeRADIUS configuration or FreeRADIUS-version-specific behavior in Go.

## Template Layout

Template sets are organized by supported FreeRADIUS version.

Example:

```text
internal/
    templates/
        3.2.10/
            clients.conf
            mods-available/
                sql
            sites-available/
```

The rendering pipeline selects the appropriate template set using the tenant's configured `RADIUSServer.Version`.

Every managed template contained within the selected version-specific template set participates in generation.

The directory structure of the template set defines the directory structure of the generated managed configuration tree.

The presence of an embedded template set for a FreeRADIUS version determines whether that version is supported by RADIUS Director.