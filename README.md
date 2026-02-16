# argocd-preview-generator

ArgoCD [Config Management Plugin](https://argo-cd.readthedocs.io/en/stable/operator-manual/config-management-plugins/) that generates preview environments from an `apps.yaml` manifest.

Reads `ARGOCD_APP_SOURCE_PATH`, loads `apps.yaml`, and outputs ArgoCD resources to stdout:
- `AppProject` (wave -3)
- `Namespace` (wave -2)
- Infra `Application`s — Helm charts (wave 0)
- Service `Application`s — pinned image tags (wave 1)

## apps.yaml

```yaml
namespace: preview-feature-xyz
services:
  - name: backend-1
    image_tag: "abc123def456"
infra:
  - name: postgres
    chart: postgresql
    repoURL: https://charts.bitnami.com/bitnami
    targetRevision: "*"
    values:
      auth:
        postgresPassword: postgres
```

## CMP sidecar

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ConfigManagementPlugin
metadata:
  name: preview-generator
spec:
  version: v1.0
  generate:
    command: [preview-generator]
  discover:
    find:
      glob: "apps.yaml"
```

## CI

GitHub Actions: test → lint → multi-arch Docker build pushed to `ghcr.io`.
