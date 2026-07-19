# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository purpose

This is a **DevOps/Cloud/Kubernetes portfolio**: a collection of small, independent, self-contained demo projects. Each one demonstrates a specific tool or technique (GitOps, IaC, Kubernetes operators, observability, admission control, etc.) in isolation. There is no shared build system, no monorepo tooling, and no cross-project dependencies — treat each leaf directory as its own project with its own toolchain.

## Repository structure

```
devops/                    CI/CD and GitOps tooling
  argo-cd/                 Argo CD bootstrapped via Helmfile; deploys NGINX through an Argo CD Application
  github-runner-kubernetes/  GitHub Actions Runner Controller (ARC) on Kubernetes via Terraform + Helm
devsecops/
  kyverno-policies/        Kyverno ClusterPolicies (disallow-privileged-containers, require-labels) enforced via admission webhook
golang/                    Go projects, each with its own go.mod (no shared module)
  http-server-kubernetes/  Go HTTP service (/version endpoint) + Dockerfile + K8s manifests (deploy/kubernetes/)
  kubernetes-api-server-ip/  client-go CLI that resolves the API server IP from EndpointSlices
  kubernetes-operator/     Kubebuilder operator; actual code lives in app-operator/ subdirectory
  unused-secret/           client-go CLI that detects Secrets unreferenced by any workload
iac/
  terraform-modules/       Reusable Terraform modules (helm_release wrappers): traefik/, ingress-nginx/
kubernetes/                Cluster bootstrap and routing demos
  eks/                     Terraform: VPC + EKS cluster + Traefik Ingress Controller
  gateway-api/             Traefik as Gateway API implementation (GatewayClass/Gateway/HTTPRoute) + manifests/
  k3d/                     Local k3d cluster + NGINX Ingress Controller (command-driven, no manifest files)
  kagent-ollama-oomkilled/ kagent (AI agent) + Ollama demo: diagnoses/remediates an OOMKilled pod
observability/
  opentelemetry/           OTel Operator auto-instrumentation demo; Collector routes OTLP to Tempo/Prometheus/Loki, visualized in Grafana; deployed via helmfile.yaml.gotmpl
  prometheus-grafana-loki/ kube-prometheus-stack + Loki/Promtail via Terraform, exposed through Traefik ingress
python/
  pulumi-eks/               Pulumi (Python) program provisioning VPC + EKS
```

Every project's `README.md` follows a consistent template: **Overview → Goals → (Architecture) → Repository Structure → Prerequisites → Usage/Deployment → Validation → Cleanup**. When adding a new project or modifying an existing one, follow this same structure — it's the convention across the whole repo, not just a suggestion for one folder.

## Working in a project

Always `cd` into the specific project directory before running any tool — there is no root-level command that operates across projects. Each project's README is the source of truth for exact commands (versions, flags, namespaces); read it before making changes rather than assuming a command works the same way in a sibling project (e.g. Helm chart versions and `--set` flags differ per project even when the same chart, like Traefik, is reused).

### Go projects (`golang/*`)
Each Go project has its own `go.mod`/`go.sum` — run `go` commands from inside the specific project directory, not the repo root.
- `go mod tidy && go run main.go` — run a CLI tool (`kubernetes-api-server-ip`, `unused-secret`)
- `golangci-lint fmt && golangci-lint run` — format/lint (used in `http-server-kubernetes`)
- `-kubeconfig=/path/to/kubeconfig` flag overrides the default `~/.kube/config` on the client-go tools

### The Kubebuilder operator (`golang/kubernetes-operator/app-operator/`)
This is the one project with a real Makefile-driven workflow (standard Kubebuilder scaffold):
- `make manifests generate` — regenerate CRDs/RBAC/DeepCopy after changing `api/v1` types
- `make fmt vet lint` — format, vet, and lint (`golangci-lint`, config in `.golangci.yml`)
- `make test` — run envtest-based tests (runs manifests/generate/fmt/vet first)
- `make test-e2e` — end-to-end tests (`test/e2e/`)
- `make run` — run the controller locally against the cluster in `~/.kube/config`
- `make install` / `make uninstall` — install/remove CRDs into the cluster
- `make docker-build docker-push IMG=<registry>/app-operator:<tag>` then `make deploy IMG=...` — build and deploy to cluster
- `make undeploy` — remove the controller deployment
- `make help` — list all targets

### Terraform projects (`devops/github-runner-kubernetes`, `kubernetes/eks`, `observability/prometheus-grafana-loki`, `iac/terraform-modules/*`)
Standard `terraform init` / `terraform plan` / `terraform apply` / `terraform destroy` per directory. Required inputs are passed as `TF_VAR_*` environment variables (e.g. `TF_VAR_github_token`, `TF_VAR_ingress_hostname`) — check each project's `variables.tf` and README for the exact set before running. `kubernetes/eks` additionally configures a remote S3 backend at `terraform init` time via `-backend-config`.

`iac/terraform-modules/{traefik,ingress-nginx}` are thin reusable `helm_release` wrapper modules (name/namespace/values inputs) meant to be consumed by other Terraform projects, not applied standalone.

### Pulumi project (`python/pulumi-eks`)
Python-based Pulumi program. Typical flow: `pulumi login --local` → `pulumi stack init` → `pulumi config set ...` (see README for the full set of required config keys, including array-valued subnet/AZ config via `--path`) → `pulumi install` → `pulumi preview` → `pulumi up`. Destroy with `pulumi destroy` then `pulumi stack rm`.

### Helm/Helmfile-driven Kubernetes projects
Several projects deploy purely via `helm`/`helmfile` + `kubectl apply` against an existing cluster (no Terraform, no app code): `devops/argo-cd`, `kubernetes/gateway-api`, `kubernetes/k3d`, `kubernetes/kagent-ollama-oomkilled`, `devsecops/kyverno-policies`, `observability/opentelemetry`. For these, the README's numbered Deployment/Usage steps *are* the build process — apply them in order, and use the Validation section's `kubectl`/`curl` commands to confirm success. Several depend on a local `k3d` cluster created with Traefik disabled (`--k3s-arg "--disable=traefik@server:*"`) since these projects install their own ingress/gateway controller.

## Cross-cutting conventions

- **HTTP-only by design, not by oversight**: many projects intentionally expose plain HTTP (port 80) for demo/validation simplicity and say so explicitly in a "Security Warning" section in their README. Don't "fix" this to HTTPS unless asked — it's a documented tradeoff for local/demo use, not a bug.
- **Namespacing**: each tool/component is installed into its own dedicated namespace (`argo`, `traefik`, `kyverno`, `kagent`, `ollama`, `oom-demo`, `loki`, `prometheus`, etc.), matching the namespace used in that project's README validation commands.
- **No manifests in `kubernetes/k3d`**: unlike other kubernetes/ projects, it's entirely command-driven (`kubectl create`/`helm install` inline) — don't look for YAML files there.
- Screenshots referenced in some READMEs live in a sibling `images/` directory within that project.
