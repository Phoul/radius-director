# ADR-0001 - Generated Configuration

**Status:** Accepted

---

# Context

Traditional FreeRADIUS deployments are maintained by manually editing configuration files.

As deployments grow, configuration becomes increasingly difficult to review, validate, and reproduce.

---

# Decision

RADIUS Director will generate FreeRADIUS configuration from a declarative domain model.

Generated configuration is considered an implementation artifact rather than the primary source of truth.

The source of truth is the configuration model maintained in version control.

---

# Consequences

Benefits include:

- reproducible deployments
- deterministic configuration
- validation before deployment
- infrastructure as code
- easier code review
- simplified automation

Administrators should avoid manually editing generated configuration.

Changes should be made to the source model and regenerated.