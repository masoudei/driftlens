#!/usr/bin/env bash
# kind-argocd-env.sh — Outputs ArgoCD env vars for DriftLens by reading
#                        them from a running Kind cluster.
#
# Usage:
#   eval "$(bash scripts/kind-argocd-env.sh)"
#   go run ./cmd/driftlens
#
# Prerequisites:
#   - kind cluster "driftlens" with ArgoCD installed
#   - kubectl, argocd CLI on PATH
#   - ArgoCD API exposed on localhost:8443
# ─────────────────────────────────────────────────────────────

CLUSTER_NAME="${CLUSTER_NAME:-driftlens}"
ARGOCD_NAMESPACE="${ARGOCD_NAMESPACE:-argocd}"
ARGOCD_URL="${ARGOCD_URL:-https://localhost:8443}"
ARGOCD_POLL_INTERVAL="${ARGOCD_POLL_INTERVAL:-30}"

# Strip protocol prefix for argocd CLI (it only accepts host:port)
ARGOCD_CLI_ADDR="${ARGOCD_URL#https://}"
ARGOCD_CLI_ADDR="${ARGOCD_CLI_ADDR#http://}"

# -- Checks ---------------------------------------------------
if ! kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
    echo "[kind-argocd-env] Kind cluster '${CLUSTER_NAME}' not found." >&2
    echo "[kind-argocd-env] Run: make dev-kind-cluster" >&2
    exit 1
fi

if ! kubectl get namespace "${ARGOCD_NAMESPACE}" &>/dev/null; then
    echo "[kind-argocd-env] ArgoCD namespace '${ARGOCD_NAMESPACE}' not found in cluster '${CLUSTER_NAME}'." >&2
    echo "[kind-argocd-env] Run: make dev-argocd-setup" >&2
    exit 1
fi

# -- Fetch admin password (or fall back to "admin") -----------
ADMIN_PWD=$(kubectl get secret argocd-initial-admin-secret \
    -n "${ARGOCD_NAMESPACE}" \
    -o jsonpath='{.data.password}' 2>/dev/null | base64 -d 2>/dev/null || echo "admin")

# -- Login and generate token ---------------------------------
TOKEN=""
if command -v argocd &>/dev/null; then
    LOGIN_OUT=$(timeout 10 argocd login "${ARGOCD_CLI_ADDR}" \
        --username admin \
        --password "${ADMIN_PWD}" \
        --insecure --grpc-web 2>&1) || true

    if echo "${LOGIN_OUT}" | grep -qi "successfully"; then
        TOKEN=$(argocd account generate-token 2>/dev/null) || true
    else
        echo "[kind-argocd-env] argocd login failed: ${LOGIN_OUT}" >&2
    fi
else
    echo "[kind-argocd-env] argocd CLI not found on PATH" >&2
fi

if [ -z "${TOKEN}" ]; then
    echo "[kind-argocd-env] Could not generate ArgoCD API token." >&2
    echo "[kind-argocd-env] ArgoCD collector will be disabled." >&2
    echo "[kind-argocd-env] (Run 'make dev-argocd-setup' first if ArgoCD is not installed)" >&2
    exit 1
fi

# -- Output env vars for eval ---------------------------------
cat <<ENV
export DRIFTLENS_ARGOCD_URL=${ARGOCD_URL}
export DRIFTLENS_ARGOCD_TOKEN=${TOKEN}
export DRIFTLENS_ARGOCD_POLL_INTERVAL=${ARGOCD_POLL_INTERVAL}
ENV
