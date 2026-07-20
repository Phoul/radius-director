# Vision

## Purpose

RADIUS Director is an open-source platform for configuring, validating, deploying, and operating standards-based RADIUS infrastructure.

Rather than manually editing FreeRADIUS configuration files, administrators describe their infrastructure using a reusable domain model. RADIUS Director validates the model and generates production-ready FreeRADIUS configuration.

The long-term goal is to make RADIUS infrastructure easier to understand, maintain, version, and deploy without replacing the proven capabilities of FreeRADIUS.

---

# Why

Large FreeRADIUS deployments often evolve over many years.

As deployments grow, configuration tends to become:

- difficult to understand
- difficult to review
- inconsistent between servers
- dependent on tribal knowledge
- risky to modify

RADIUS Director aims to solve these operational problems while preserving the flexibility and performance of FreeRADIUS.

---

# Design Principles

## Standards First

The project models standards, not products.

Vendor-specific implementations should be represented only where necessary.

Examples include:

- RFC 2865 (Authentication)
- RFC 2866 (Accounting)
- RFC 5176 (Disconnect and CoA)

Applications such as billing systems or management platforms interact with RADIUS Director through these standards rather than through product-specific integrations.

---

## FreeRADIUS Remains the RADIUS Server

RADIUS Director is not a replacement for FreeRADIUS.

FreeRADIUS remains responsible for:

- Authentication
- Authorization
- Accounting
- Proxying
- Policy execution

RADIUS Director generates and manages its configuration.

---

## Infrastructure as Code

Configuration should be declarative.

Infrastructure should be reproducible.

Every deployment should be generated from version-controlled configuration.

---

## Validation Before Deployment

Configuration should be validated before deployment.

Configuration errors should be detected before they reach production.

---

## Vendor Neutral

The project should not assume a specific NAS vendor, billing platform, or monitoring system.

---

## Human Readable

Generated configuration should remain understandable by experienced FreeRADIUS administrators.

---

## Reusable Global Objects

Reusable configuration should be defined once and referenced wherever possible.

Tenant-specific infrastructure should remain isolated.

This minimizes duplication while allowing multiple tenants to share common operational definitions.

---

# Non-Goals

RADIUS Director is not intended to:

- replace FreeRADIUS
- replace a billing platform
- become a captive portal
- become a full network management platform
- require a web interface

The command line and generated configuration are first-class interfaces.