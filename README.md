# DriftLens

**Infrastructure Change Intelligence.**

DriftLens is a GitOps Intelligence Platform that transforms infrastructure drift from an alert into an explainable event. It builds a causality graph connecting Git commits, CI/CD pipelines, Terraform/OpenTofu runs, ArgoCD/Flux events, and Kubernetes audit logs to answer:

- Who changed production?
- What changed?
- Why did it change?
- What was affected?
- How risky is it?
- How do we fix it?

## Architecture

```
Collector → Event Bus → Correlator → Graph Engine → Risk Engine → API → UI
```

### Core Components

| Component | Description |
|---|---|
| **Event Collectors** | Ingest from K8s audit logs, ArgoCD, Flux, Terraform, OpenTofu, GitHub, GitLab, OpenTelemetry |
| **Correlation Engine** | Builds causality relationships between events |
| **Graph Engine** | Stores and queries the causality graph (nodes + relationships) |
| **Drift Engine** | Detects configuration, runtime, and infrastructure drift |
| **Attribution Engine** | Identifies the user, service account, or CI job responsible |
| **Risk Engine** | Calculates blast radius, criticality, security & compliance impact |
| **AI Analysis Engine** | Produces root cause analysis, remediation suggestions, drift narratives |

## Graph Model

Everything in DriftLens is a node connected by relationships:

```
Commit ──CAUSED──► Pipeline
Pipeline ──TRIGGERED──► TerraformRun
TerraformRun ──DEPLOYED──► ArgoSync
ArgoSync ──MODIFIED──► Deployment
Deployment ──DRIFTED──► DriftEvent
DriftEvent ──IMPACTS──► Service
```

### Node Types

- `Commit` — Git commit
- `Pipeline` — CI/CD pipeline run
- `TerraformRun` — Terraform/OpenTofu apply
- `ArgoSync` — ArgoCD/Flux sync event
- `KubernetesEvent` — K8s audit event
- `DriftEvent` — Detected drift occurrence

### Relationship Types

- `CAUSED` — direct causation
- `TRIGGERED` — pipeline or automation trigger
- `MODIFIED` — resource modification
- `DEPLOYED` — deployment action
- `DRIFTED` — drift detected
- `IMPACTS` — service impact

## Getting Started

```bash
git clone https://github.com/masoudei/driftlens.git
cd driftlens
go build ./cmd/driftlens
./driftlens
```

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go, Gin |
| Database | PostgreSQL, OpenSearch |
| Event Bus | NATS |
| Frontend | Next.js, TypeScript, Tailwind |
| Infrastructure | Kubernetes, Helm |
| Observability | OpenTelemetry |

## Roadmap

| Phase | Focus |
|---|---|
| **Phase 1** | MVP — K8s audit collector, ArgoCD integration, drift detection, event timeline, user attribution |
| **Phase 2** | Multi-cluster, Terraform/GitHub/GitLab correlation |
| **Phase 3** | OpenTelemetry correlation, risk scoring, blast radius, historical analytics |
| **Phase 4** | AI — root cause analysis, drift summaries, incident reconstruction, remediation |

## Contributing

1. Branch off `main`: `git checkout -b feat/my-feature`
2. Write code following the [design principles](docs/design-principles.md)
3. Add tests
4. Open a pull request

See [AGENTS.md](AGENTS.md) for detailed development workflow and conventions.

## License

MIT
