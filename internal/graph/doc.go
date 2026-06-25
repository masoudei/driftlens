// Package graph implements the core causality graph for DriftLens.
//
// DriftLens uses graph thinking: everything is a node connected by relationships.
//
// Node Types:
//   - Commit
//   - Pipeline
//   - TerraformRun
//   - ArgoSync
//   - KubernetesEvent
//   - DriftEvent
//
// Relationship Types:
//   - CAUSED
//   - TRIGGERED
//   - MODIFIED
//   - DEPLOYED
//   - DRIFTED
//   - IMPACTS
//
// Example:
//
//	Commit ──CAUSED──► Pipeline ──TRIGGERED──► TerraformRun
//	TerraformRun ──DEPLOYED──► ArgoSync ──MODIFIED──► Deployment
//	Deployment ──DRIFTED──► DriftEvent ──IMPACTS──► Service
package graph
