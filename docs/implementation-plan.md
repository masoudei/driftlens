# Implementation Prompt

Build DriftLens.

A cloud-native GitOps Intelligence Platform.

Technology stack:

Backend:
- Go
- Gin
- PostgreSQL
- OpenSearch
- NATS

Frontend:
- Next.js
- TypeScript
- Tailwind

Infrastructure:
- Kubernetes
- Helm
- OpenTelemetry

Requirements:

## MVP

Implement:

1. Kubernetes Audit Collector

Watch:

- Deployment
- StatefulSet
- DaemonSet
- ConfigMap
- Secret

changes.

Store events.

---

2. ArgoCD Integration

Collect:

- Sync Events
- Health Events
- Drift Events

Store relationships.

---

3. Event Timeline

Generate timelines:

09:00 Git Commit
09:05 Pipeline
09:10 Terraform
09:20 ArgoCD
09:25 Drift

---

4. Attribution Engine

Identify:

- User
- ServiceAccount
- CI Job

responsible for change.

---

5. REST API

Endpoints:

GET /drifts

GET /drifts/{id}

GET /timeline/{resource}

GET /risk/{resource}

---

6. Frontend

Pages:

Dashboard

Drift Explorer

Timeline View

Risk Analysis

Resource Explorer

---

7. Future Ready

Design graph model supporting:

Node Types:

- Commit
- Pipeline
- TerraformRun
- ArgoSync
- KubernetesEvent
- DriftEvent

Relationship Types:

- CAUSED
- TRIGGERED
- MODIFIED
- DEPLOYED
- DRIFTED

Graph queries should support root-cause analysis.

The platform should be CNCF-style and scalable to 100+ clusters.