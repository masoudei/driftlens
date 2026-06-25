#!/usr/bin/env bash
set -euo pipefail

# ─────────────────────────────────────────────────────────────
# setup-argocd-kind.sh
# ─────────────────────────────────────────────────────────────
# Installs ArgoCD on a Kind cluster, creates a sample app,
# generates an API token, and prints env vars for DriftLens.
#
# Prerequisites:
#   - kind, kubectl, argocd CLI
#   - kind cluster named "driftlens" (run: kind create cluster --config kind-config.yaml --name driftlens)
#
# Usage:
#   bash scripts/setup-argocd-kind.sh
# ─────────────────────────────────────────────────────────────

CLUSTER_NAME="${CLUSTER_NAME:-driftlens}"
ARGOCD_VERSION="${ARGOCD_VERSION:-v2.14.0}"
NAMESPACE="${NAMESPACE:-argocd}"

echo "═══ ArgoCD Setup for DriftLens ═══"
echo "  Cluster : $CLUSTER_NAME"
echo "  Version : $ARGOCD_VERSION"
echo ""

# ── 1. Ensure cluster exists ────────────────────────────────
echo "▸ Checking Kind cluster …"
if ! kind get clusters | grep -q "^${CLUSTER_NAME}$"; then
  echo "  Cluster '${CLUSTER_NAME}' not found. Creating …"
  kind create cluster --config kind-config.yaml --name "${CLUSTER_NAME}"
fi
kubectl cluster-info --context "kind-${CLUSTER_NAME}" > /dev/null 2>&1
echo "  OK"

# ── 2. Install ArgoCD ───────────────────────────────────────
echo "▸ Installing ArgoCD ${ARGOCD_VERSION} …"
kubectl create namespace "${NAMESPACE}" --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -n "${NAMESPACE}" -f "https://raw.githubusercontent.com/argoproj/argo-cd/${ARGOCD_VERSION}/manifests/install.yaml"

echo "  Waiting for ArgoCD components …"
kubectl wait --for=condition=Available -n "${NAMESPACE}" deployment/argocd-server --timeout=180s
kubectl wait --for=condition=Available -n "${NAMESPACE}" deployment/argocd-repo-server --timeout=180s
kubectl wait --for=condition=Available -n "${NAMESPACE}" deployment/argocd-redis --timeout=120s
echo "  OK"

# ── 3. Patch argocd-server to NodePort ──────────────────────
echo "▸ Patching argocd-server service to NodePort …"
kubectl patch svc argocd-server -n "${NAMESPACE}" -p '{"spec":{"type":"NodePort","ports":[{"name":"https","port":443,"targetPort":8080,"nodePort":30443}]}}'
echo "  OK (ArgoCD API available at https://localhost:8443)"

# ── 4. Get admin password ───────────────────────────────────
echo "▸ Retrieving admin password …"
ADMIN_PASSWORD=$(kubectl get secret argocd-initial-admin-secret -n "${NAMESPACE}" -o jsonpath="{.data.password}" 2>/dev/null | base64 -d || true)
if [ -z "$ADMIN_PASSWORD" ]; then
  # Fallback: restart server to regenerate
  echo "  Initial secret not found, restarting argocd-server to regenerate …"
  kubectl delete pod -n "${NAMESPACE}" -l app.kubernetes.io/name=argocd-server --force --grace-period=0 2>/dev/null || true
  sleep 10
  ADMIN_PASSWORD=$(kubectl get secret argocd-initial-admin-secret -n "${NAMESPACE}" -o jsonpath="{.data.password}" 2>/dev/null | base64 -d || true)
  if [ -z "$ADMIN_PASSWORD" ]; then
    ADMIN_PASSWORD="admin"
    echo "  ⚠ Could not retrieve secret, defaulting to 'admin'"
  fi
fi
echo "  Admin password: ${ADMIN_PASSWORD}"

# ── 5. Log in via argocd CLI ────────────────────────────────
echo "▸ Logging in to ArgoCD …"
argocd login localhost:8443 --username admin --password "${ADMIN_PASSWORD}" --insecure --grpc-web
echo "  OK"

# ── 6. Create a sample Application ──────────────────────────
echo "▸ Creating sample namespace 'sample-app' …"
kubectl create namespace sample-app --dry-run=client -o yaml | kubectl apply -f -

echo "▸ Registering sample-app namespace with ArgoCD …"
kubectl label namespace sample-app argocd.argoproj.io/managed-by="${NAMESPACE}" --overwrite

echo "▸ Creating ArgoCD Application 'sample-nginx' …"
argocd app create sample-nginx \
  --repo https://github.com/argoproj/argocd-example-apps.git \
  --path guestbook \
  --dest-server https://kubernetes.default.svc \
  --dest-namespace sample-app \
  --sync-policy automated \
  --auto-prune \
  --self-heal

echo "  Waiting for initial sync …"
sleep 10
argocd app sync sample-nginx
echo "  OK"

# ── 7. Create API token for DriftLens ───────────────────────
echo "▸ Creating DriftLens API token …"
DRIFTLENS_TOKEN=$(argocd account generate-token --account admin 2>/dev/null || argocd account generate-token 2>/dev/null)
if [ -z "$DRIFTLENS_TOKEN" ]; then
  echo "  ⚠ Token generation failed. Using admin password as fallback."
  DRIFTLENS_TOKEN="${ADMIN_PASSWORD}"
fi
echo "  Token: ${DRIFTLENS_TOKEN}"

# ── 8. Print DriftLens env vars ─────────────────────────────
echo ""
echo "══════════════════════════════════════════════════════════"
echo "  DriftLens ArgoCD Connection Details"
echo "══════════════════════════════════════════════════════════"
echo ""
echo "  export DRIFTLENS_ARGOCD_URL=https://localhost:8443"
echo "  export DRIFTLENS_ARGOCD_TOKEN=${DRIFTLENS_TOKEN}"
echo "  export DRIFTLENS_ARGOCD_POLL_INTERVAL=30"
echo ""
echo "  Then run DriftLens with:"
echo "    make dev-k8s"
echo "  or"
echo "    DRIFTLENS_DATABASE_URL=postgres://driftlens:driftlens@localhost:5433/driftlens?sslmode=disable \\"
echo "      DRIFTLENS_ARGOCD_URL=https://localhost:8443 \\"
echo "      DRIFTLENS_ARGOCD_TOKEN=${DRIFTLENS_TOKEN} \\"
echo "      go run ./cmd/driftlens"
echo ""
echo "══════════════════════════════════════════════════════════"
