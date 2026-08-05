# ADR-0007: Generic Managed Template Rendering

## Status

Accepted

## Context

ADR-0004 establishes that RADIUS Director generates managed configuration from version-specific FreeRADIUS template sets.

The initial implementation rendered each managed configuration file through dedicated renderer functions. Examples included:

- RenderClients()
- RenderProxy()
- RenderSQL()
- RenderCOA()
- RenderAuthorize()

Each newly managed configuration file required corresponding changes to renderer logic, output generation, and associated tests.

As the number of managed configuration files grows, this approach unnecessarily couples renderer code to individual FreeRADIUS configuration file names and increases maintenance effort.

Future FreeRADIUS releases may add, remove, rename, or reorganize managed configuration files.

The rendering architecture should allow the managed template set to evolve without requiring renderer changes solely because managed templates have been added or removed.

---

## Decision

The renderer SHALL treat the selected version-specific template set as the authoritative definition of the managed configuration tree.

The renderer SHALL:

- select the template set associated with the configured FreeRADIUS version
- discover every managed template contained within the selected version-specific template set
- preserve the relative directory structure of each managed template
- execute each template using the generated tenant model
- produce the corresponding managed configuration tree

The renderer SHALL NOT contain knowledge of individual managed FreeRADIUS configuration file names.

Adding a new managed configuration file SHOULD normally require only adding the corresponding template to the appropriate version-specific template set.

All managed templates SHALL be executed with generator.Tenant as the root template context. Templates may reference any information reachable from the tenant object.

Every file within a managed template set SHALL be treated as a Go template. Files that contain no template actions are rendered unchanged.

---

## Responsibilities

The generator is responsible for producing the complete tenant model used during rendering.

The renderer is responsible for:

- discovering managed templates
- executing managed templates
- preserving template directory structure
- producing the managed configuration tree

The template set is responsible for defining:

- which managed configuration files exist
- where they are located within the managed configuration tree
- how each managed configuration file is rendered

---

## Consequences

Advantages:

- Eliminates renderer boilerplate for individual managed configuration files.
- Removes coupling between renderer code and specific FreeRADIUS configuration file names.
- Simplifies support for future FreeRADIUS releases.
- Allows managed configuration files to be added without modifying renderer logic.
- Makes the version-specific template set the complete definition of the managed configuration tree.
- Produces a smaller and simpler renderer implementation.

Trade-offs:

- The renderer becomes dependent upon the structure of the version-specific template set.
- Template organization becomes part of the project's public architecture.

These trade-offs are preferable to maintaining renderer logic for every managed configuration file.

---

## Relationship to ADR-0004

ADR-0004 defines the purpose, responsibilities, and organization of managed FreeRADIUS templates.

This ADR defines how those managed templates are discovered and rendered.

Together, the two ADRs establish that:

- the generator produces the tenant model
- the version-specific template set defines the managed configuration tree
- the renderer generically renders every managed template contained within the selected template set