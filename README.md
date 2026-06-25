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

### Prerequisites

- Go 1.25+
- Kubernetes cluster (optional — dev mode works without one)
- Kind (optional — for local K8s dev)

### Clone & Build

```bash
git clone https://github.com/masoudei/driftlens.git
cd driftlens
go build ./cmd/driftlens
```

## Usage

### One command (PostgreSQL + app)

```bash
make dev
```

Auto-starts PostgreSQL via Docker Compose if not running, then launches DriftLens with mock collector + persistent storage. Idempotent — safe to run repeatedly.

### Other modes

```bash
make dev-db-only          # start PostgreSQL only
make dev-kind             # run against Kind cluster (in-memory)
make dev-kind-db          # run against Kind + PostgreSQL
```

### With Kind (local K8s cluster)

```bash
make dev-kind-cluster     # create the cluster
make dev-kind             # build and connect to the cluster
```

Kind automatically merges the cluster context into your default kubeconfig (`~/.kube/config`). No need to set `KUBECONFIG`.

To verify the context was added:

```bash
kubectl cluster-info --context kind-driftlens
```

To get the raw kubeconfig (e.g. for CI or other tools):

```bash
kind get kubeconfig --name driftlens > driftlens-kubeconfig.yaml
KUBECONFIG=driftlens-kubeconfig.yaml go run ./cmd/driftlens
```

### Makefile targets

```bash
make dev              # dev mode (mock collector, in-memory store)
make dev-db           # start PostgreSQL
make dev-db-run       # dev mode with PostgreSQL
make test             # run all tests
make test-cover       # tests + coverage report
make build            # compile binary
make clean            # remove build artifacts
```

### Dummy Resources

Create sample resources and modify them to trigger drift events:

```bash
# Create resources
kubectl create deployment demo-nginx --image=nginx --replicas=2
kubectl create configmap demo-config --from-literal=color=blue --from-literal=env=staging
kubectl create secret generic demo-secret --from-literal=password=secret123
kubectl create namespace demo-env

# Scale deployment (triggers replicas drift)
kubectl scale deployment/demo-nginx --replicas=5

# Update configmap (triggers data_keys drift)
kubectl create configmap demo-config --from-literal=color=red --from-literal=env=prod --from-literal=region=us-east -o yaml --dry-run=client | kubectl apply -f -

# Update secret (triggers secret_keys drift)
kubectl create secret generic demo-secret --from-literal=password=newpass456 --from-literal=token=abc123 -o yaml --dry-run=client | kubectl apply -f -

# Delete resources
kubectl delete deployment demo-nginx
kubectl delete configmap demo-config
kubectl delete secret demo-secret
kubectl delete namespace demo-env
```

### ArgoCD Integration

```bash
# 1. Ensure Kind cluster running with port 8443 mapped
make dev-kind-cluster

# 2. Install ArgoCD, create sample app, generate token
bash scripts/setup-argocd-kind.sh

# 3. Set env vars (printed by the script)
set DRIFTLENS_ARGOCD_URL=https://localhost:8443
set DRIFTLENS_ARGOCD_TOKEN=<token-from-script>

# 4. Run DriftLens connected to Kind + ArgoCD
make dev-kind
```

The script:
- Installs ArgoCD v2.14.0 on the Kind cluster
- Patches argocd-server to NodePort (30443 → localhost:8443)
- Creates a sample Guestbook app synced from `argocd-example-apps`
- Generates an API token for DriftLens
- Prints the required env vars

To clean up:
```bash
make dev-argocd-clean
```

### Query Endpoints

```bash
# All detected drifts
curl localhost:8081/drifts

# Single drift detail
curl localhost:8081/drifts/drift-<id>

# Event timeline for a resource type
curl localhost:8081/timeline/Deployment
curl localhost:8081/timeline/ArgoSync

# Risk assessment for a resource type
curl localhost:8081/risk/Deployment
```

### Frontend

```bash
make web     # start dev server at http://localhost:3000
```

The frontend proxies API calls through Next.js — no CORS needed. Requires the backend to be running on port 8081.

### Environment variables

| Variable | Default | Description |
|---|---|---|
| `DRIFTLENS_DEV` | `false` | Set to `true` to use mock collector (no K8s) |
| `DRIFTLENS_DATABASE_URL` | — | PostgreSQL connection string (empty = in-memory store) |
| `DRIFTLENS_NAMESPACE` | all | K8s namespace to watch (empty = all) |
| `DRIFTLENS_ADDR` | `:8081` | API server address |
| `DRIFTLENS_ARGOCD_URL` | — | ArgoCD server URL (optional, enables ArgoCD collector) |
| `DRIFTLENS_ARGOCD_TOKEN` | — | ArgoCD auth token |
| `DRIFTLENS_ARGOCD_POLL_INTERVAL` | `30` | ArgoCD poll interval in seconds |
| `KUBECONFIG` | `~/.kube/config` | K8s config file path |

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
