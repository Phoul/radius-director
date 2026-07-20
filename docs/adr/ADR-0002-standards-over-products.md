# ADR-0002 - Standards Over Products

**Status:** Accepted

---

# Context

RADIUS deployments frequently integrate with many third-party products including:

- Billing platforms
- Network management systems
- Monitoring platforms
- Identity providers
- Network vendors

These products evolve independently and may be replaced over time.

If the internal data model is built around specific products, the architecture becomes tightly coupled to those implementations and difficult to adapt.

---

# Decision

RADIUS Director models industry standards and operational concepts rather than individual products.

The domain model should represent concepts defined by RADIUS and related standards whenever possible.

Examples include:

- Authentication
- Authorization
- Accounting
- Disconnect Requests
- Change of Authorization (CoA)
- RADIUS Clients
- RADIUS Servers
- Shared Secrets
- Session Accounting

Vendor-specific behaviour should only be introduced where it represents a capability that cannot be expressed through a standard model.

---

# Examples

## Preferred

Model:

- Credential Profile
- Authentication Profile
- Accounting Profile
- NAS
- RADIUS Server
- Database
- Deployment

These concepts remain valid regardless of the products in use.

---

## Avoid

Avoid creating objects such as:

- Sonar Server
- MikroTik Configuration
- Cisco Configuration
- Ubiquiti Profile
- PowerDNS Integration

These are implementations of broader concepts rather than concepts themselves.

---

# Product Integrations

External products should interact with RADIUS Director through standards whenever possible.

Examples include:

- Billing systems sending RFC 5176 Disconnect Requests
- NAS devices communicating using standard RADIUS authentication and accounting
- Monitoring systems performing authentication or accounting tests
- Databases accessed through supported SQL drivers

RADIUS Director should not require knowledge of the business purpose of these systems.

---

# Exceptions

Some vendor-specific behaviour is unavoidable.

Examples include:

- Vendor Specific Attributes (VSAs)
- Vendor-specific SQL queries
- Device-specific authentication policies
- Vendor-specific monitoring capabilities

These should be isolated and clearly identified within the configuration model.

---

# Consequences

Benefits include:

- Vendor-neutral architecture
- Easier migration between platforms
- Longer-term maintainability
- Reduced coupling
- Simpler documentation
- Easier testing

Trade-offs include:

- Some vendor-specific features may require extensions to the domain model.
- Certain advanced features may not have a completely vendor-neutral representation.

---

# Rationale

Products change.

Standards evolve much more slowly.

By modelling standards rather than products, RADIUS Director becomes resilient to changes in vendors, management platforms, and operational workflows.

The architecture remains focused on describing a RADIUS deployment rather than the software currently interacting with it.