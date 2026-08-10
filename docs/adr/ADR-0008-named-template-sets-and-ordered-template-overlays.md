# ADR-0008: Named Template Sets and Ordered Template Overlays

- Status: Accepted
- Date: 2026-08-10

## Context

RADIUS Director manages FreeRADIUS configuration as version-specific template trees.

The original design selected the managed template tree primarily from the FreeRADIUS version. This worked while each FreeRADIUS version had a single configuration variant, but it does not provide a clean mechanism for supporting multiple configuration variants for the same FreeRADIUS version or for making a small, deployment-specific modification without duplicating an entire template tree.

RADIUS Director now supports deployment profiles. A deployment profile selects a named base template set and may additionally specify an ordered list of overlays.

This establishes a distinction between:

- the FreeRADIUS version being deployed;
- the catalog of templates available for that version;
- the base template set selected for a deployment;
- optional overlays applied to that base template set.

The resulting configuration is therefore not necessarily the contents of one template directory. It is an effective template tree composed from a base template set and zero or more ordered overlays.

## Decision

### 1. Template catalogs are scoped by FreeRADIUS version

Template sets are associated with a specific FreeRADIUS version.

A version determines which template sets and overlays are available for that version. The version does not, by itself, select the template set used for a tenant.

Conceptually:

    FreeRADIUS version
            |
            v
    version-scoped template catalog
            |
            +-- named base template sets
            |
            +-- named overlays

A template set or overlay intended for one FreeRADIUS version is not implicitly available for another version.

### 2. Deployment profiles select the base template set

A deployment profile contains a `template` identifier selecting the named base template set to use.

For example:

    deployment_profiles:
      default:
        template: default

The template identifier is a name, not a filesystem path.

The selected template set provides the initial configuration tree for the tenant.

### 3. Deployment profiles may specify ordered overlays

A deployment profile may contain an ordered `overlays` list.

For example:

    deployment_profiles:
      experimental:
        template: default
        overlays:
          - coa-relay-test
          - debug-logging

The order of the list is significant.

Overlays are applied in the order in which they are specified.

### 4. Overlay precedence is last-wins

Each overlay contains files using paths relative to the effective template root.

When an overlay contains a file with the same effective path as a file already present in the base template or an earlier overlay, the later overlay replaces the earlier file.

For example:

    base:
      sites-available/coa

    overlay-a:
      sites-available/coa

    overlay-b:
      sites-available/coa

with:

    overlays:
      - overlay-a
      - overlay-b

results in:

    sites-available/coa
        -> overlay-b version

This same rule applies regardless of whether the replaced file originated in the base template set or an earlier overlay.

### 5. Overlays may add files

An overlay may contain files that do not exist in the base template set.

Those files become part of the effective template tree.

For example:

    base:
      sites-available/default

    overlay:
      sites-available/example

results in an effective tree containing both files.

### 6. Overlays do not delete files

Overlays are additive/replacement-only.

The absence of a file from an overlay does not remove a file from the base template set or an earlier overlay.

RADIUS Director does not currently provide a deletion or tombstone mechanism for overlays.

If a future requirement requires removing a file from a base template set, that should be addressed through an explicit future design rather than by assigning special meaning to file absence.

### 7. Template and overlay identifiers are names

Template-set and overlay identifiers represent named objects, not filesystem paths.

They must therefore be valid single path components rather than arbitrary relative paths.

Path separators are not permitted, and `.` is not a valid identifier.

Paths within a template or overlay may contain directories. For example:

    sites-available/coa

is a valid managed file path, while:

    experimental/coa

is not a valid overlay identifier.

### 8. The effective template tree is resolved before rendering

Template resolution produces an effective mapping of managed relative paths to the physical template files that supply them.

The renderer operates on the effective tree and does not need to know whether an individual file came from:

- the base template set;
- the first overlay;
- a later overlay.

This keeps overlay composition inside the template subsystem rather than distributing overlay-specific logic throughout rendering.

### 9. Validation must resolve the selected template configuration

A deployment profile is not considered fully valid merely because its `template` and `overlays` fields are syntactically valid.

The selected base template set and each referenced overlay must be available for the tenant's FreeRADIUS version.

Validation should identify these failures before generation/rendering rather than allowing an otherwise valid configuration to fail later during rendering.

Where practical, independent template-resolution errors should be reported together with other validation errors.

## Consequences

### Positive consequences

- Multiple configuration variants can coexist for the same FreeRADIUS version.
- Small experimental or customer-specific changes do not require duplicating an entire template tree.
- Overlay ordering provides deterministic precedence.
- A later overlay can replace a file without modifying the base template set.
- New managed files can be introduced by an overlay.
- The renderer remains unaware of the source of individual effective files.
- Template selection becomes an explicit deployment-profile decision rather than an implicit consequence of the FreeRADIUS version.

### Negative consequences

- The effective configuration is more complex than a single template directory.
- Template resolution must account for the FreeRADIUS version, base template set, and ordered overlays.
- Validation must verify that the selected template resources actually exist.
- Multiple overlays can make it harder to determine which file is ultimately effective if their use is not documented clearly.
- Testing must cover both individual overlays and overlay precedence.

## Relationship to Previous Decisions

This ADR supersedes the portions of ADR-0004 and ADR-0007 that define the FreeRADIUS version as selecting a single complete template tree.

Those decisions remain applicable where they do not conflict with this ADR.

In particular, the managed template concept and version-specific nature of the FreeRADIUS configuration remain valid; the difference is that the effective tree may now be composed from a named base template set and ordered overlays.

ADR-0005 remains a separate concern regarding ownership and materialization of the generated configuration tree.

## Example

A deployment profile may be defined as:

    deployment_profiles:
      default:
        template: default
        overlays: []

      coa-relay-test:
        template: default
        overlays:
          - coa-relay-test

For a tenant using `coa-relay-test` with FreeRADIUS 3.2.10, RADIUS Director resolves:

    3.2.10/default/
        +
    3.2.10/overlays/coa-relay-test/

into one effective template tree.

If both trees contain:

    sites-available/coa

the overlay version is used.

If the overlay contains a new file such as:

    sites-available/coa-relay

that file is added to the effective tree.