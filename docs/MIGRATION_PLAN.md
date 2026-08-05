# Production Migration Plan

This document tracks the migration from the existing production FreeRADIUS deployment to RADIUS Director.

---

# Phase 0 - Documentation

- [x] Define project vision
- [x] Define architecture
- [ ] Finalize domain model
- [ ] Finalize object model
- [ ] Finalize configuration schema
- [ ] Record architectural decisions

---

# Phase 1 - Discovery

- [x] Reverse engineer existing deployment
- [x] Document CoA architecture
- [x] Document authentication flow
- [x] Document accounting flow
- [ ] Document all generated FreeRADIUS files

---

# Phase 2 - Design

- [ ] Finalize architecture
- [ ] Finalize domain model
- [ ] Finalize object model
- [ ] Finalize configuration schema
- [ ] Define generator architecture

---

# Phase 3 - Managed Configuration Generation

- [x] Parse configuration
- [x] Resolve object relationships
- [x] Validate configuration
- [x] Generate clients.conf
- [ ] Generate radiusd.conf
- [x] Generate proxy.conf
- [ ] Generate SQL configuration
- [ ] Generate managed virtual server configuration
- [x] Assemble managed configuration tree

---

# Phase 4 - Deployment

- [ ] Materialize managed configuration trees on the target filesystem
- [ ] Deploy managed configuration
- [ ] Start FreeRADIUS

---

# Phase 5 - Testing

- [ ] Deploy lab environment
- [ ] Verify authentication
- [ ] Verify accounting
- [ ] Verify CoA
- [ ] Validate generated configuration
- [ ] Compare generated configuration with production

---

# Phase 6 - Cloud Deployment

- [ ] Build production deployment artifacts
- [ ] Deploy cloud infrastructure
- [ ] Configure databases
- [ ] Configure monitoring
- [ ] Perform production cutover

---

# Success Criteria

Migration is complete when:

- Declarative configuration replaces manually maintained FreeRADIUS configuration.
- Generated configuration is functionally equivalent to the existing production deployment.
- Production authentication is fully operational.
- Production accounting is fully operational.
- Disconnect and CoA requests function correctly.
- Operational procedures are fully documented.