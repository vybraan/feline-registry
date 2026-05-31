# Feline registry

Tiny Go API. POST /names to add, GET /names to list. I chose k3d for this.

## Start

```bash
k3d cluster create feline --agents 2
kubectl apply -f k8s/deploy/ -f k8s/rbac/
kubectl rollout status deploy/feline-registry

kubectl port-forward service/feline-registry 8080:80
curl -X POST localhost:8080/names -d '{"name":"lion"}'
curl localhost:8080/names
curl localhost:8080/metrics

kubectl auth can-i delete pod \
  --as=system:serviceaccount:ci:ci-user -n ci

ansible-galaxy collection install community.general
ansible-playbook ansible/playbook.yaml -i localhost,
```

## Structure

- `app/` - list Go source, Dockerfile
- `k8s/deploy/` - 2-replica Deployment with liveness/readiness probes
- `k8s/rbac/` - ServiceAccount + Role + RoleBinding (read-only pods/deployments in `ci` namespace)
- `ansible/` - playbook: installs deps, configures ufw, disables root SSH
- `.github/workflows/ci.yml` _ golangci-lint, gosec, tests, build + Trivy scan

Image is built & pushed to `ghcr.io/vybraan/feline-registry`.
