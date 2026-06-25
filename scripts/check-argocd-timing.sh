#!/usr/bin/env bash
set -euo pipefail

TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcmdvY2QiLCJzdWIiOiJhZG1pbjphcGlLZXkiLCJuYmYiOjE3ODIzNTE2MzQsImlhdCI6MTc4MjM1MTYzNCwianRpIjoiYjVhODVmYmEtZGY3Ny00NzE0LTk3ZmYtNDVkM2Y3NjBlNzUzIn0.xAS1ukhLgu5Myn7LyEWUjqR4p2M-WbwvLqZqYxQBP8c"

kubectl scale deployment/guestbook-ui -n sample-app --replicas=3
echo "Scaled at $(date +%H:%M:%S)"

for i in $(seq 1 12); do
  status=$(curl.exe -s -k -H "Authorization: Bearer $TOKEN" https://localhost:8443/api/v1/applications/sample-nginx 2>/dev/null | python3 -c "import sys,json;d=json.load(sys.stdin);print(d['status']['sync']['status'])" 2>/dev/null)
  echo "$(date +%H:%M:%S) $status"
  sleep 5
done
