# Architecture

## Core Components

### Event Collectors

Responsible for ingesting:

- Kubernetes Audit Logs
- ArgoCD Events
- Flux Events
- Terraform State Changes
- OpenTofu State Changes
- GitHub Events
- GitLab Events
- OpenTelemetry Traces

---

### Correlation Engine

Responsible for building relationships between events.

Example:

Commit
→ Pipeline
→ Terraform Apply
→ Argo Sync
→ Deployment Change
→ Audit Event

---

### Drift Engine

Detects:

- Configuration Drift
- Runtime Drift
- Infrastructure Drift

---

### Attribution Engine

Determines:

- User
- Service Account
- CI Job
- Automation System

that caused the drift.

---

### Risk Engine

Calculates:

- Blast Radius
- Criticality
- Security Impact
- Compliance Impact

---

### AI Analysis Engine

Produces:

- Root Cause Analysis
- Suggested Remediation
- Incident Summaries
- Drift Narratives