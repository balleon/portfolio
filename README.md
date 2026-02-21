# Portfolio

I build and operate reliable cloud-native platforms with a strong focus on **Kubernetes**, **automation**, and **operational excellence**.  
This portfolio contains hands-on projects that demonstrate practical skills in platform engineering, infrastructure as code, CI/CD, observability, and policy-driven security.

I am targeting roles as:

- **Site Reliability Engineer (SRE)**
- **Platform Engineer**
- **DevOps Engineer**

---

## Skills & Tools

**Cloud & Infrastructure:** AWS, VPC, EKS, Network Load Balancer  
**Containers & Orchestration:** Docker, Kubernetes, k3d, NGINX Ingress  
**IaC & GitOps:** Terraform, Helm, Helmfile, Argo CD  
**Observability:** Prometheus, Alertmanager, Grafana, Loki, Promtail  
**CI/CD & Automation:** GitHub Actions, self-hosted runners, scripting  
**Programming:** Go, Python

---

## Featured Projects

| Project | What it demonstrates | Key technologies |
|---|---|---|
| [GitOps with Argo CD](./projects/devops/argo-cd/) | Delivers repeatable application deployments through declarative GitOps workflows. | Argo CD, Helmfile, Kubernetes |
| [GitHub Runner on Kubernetes](./projects/devops/github-runner-kubernetes/) | Enables scalable self-hosted CI execution on Kubernetes for delivery pipelines. | Terraform, GitHub Actions Runner Controller, Helm, Kubernetes |
| [Kubernetes Policy Enforcement](./projects/devsecops/kyverno-policies/) | Enforces security and governance guardrails with Kubernetes-native policy-as-code. | Kyverno, Kubernetes |
| [HTTP Server on Kubernetes (Go)](./projects/golang/http-server-kubernetes/) | Deploys a containerized Go service with production-style Kubernetes networking. | Go, Docker, Kubernetes, Ingress |
| [Unused Secret Detector (Go)](./projects/golang/unused-secret/) | Improves cluster security posture by detecting and reporting unused Kubernetes Secrets. | Go, client-go, Kubernetes API |
| [Terraform Module: Ingress NGINX](./projects/iac/terraform-modules/ingress-nginx/) | Standardizes ingress provisioning with a reusable Terraform module pattern. | Terraform, Helm, Kubernetes |
| [AWS EKS with Terraform](./projects/kubernetes/eks/) | Provisions an AWS EKS platform foundation using infrastructure as code. | Terraform, AWS EKS, Helm, NGINX Ingress |
| [Monitoring & Observability](./projects/observability/prometheus-grafana-loki/) | Implements a full observability stack for metrics, logs, dashboards, and alerting. | Prometheus, Alertmanager, Grafana, Loki, Promtail, Terraform |

---

## Projects by Section

| Section | Projects |
|---|---|
| **DevOps** | [Argo CD GitOps](./projects/devops/argo-cd/), [GitHub Runner Kubernetes](./projects/devops/github-runner-kubernetes/) |
| **DevSecOps** | [Kyverno Policies](./projects/devsecops/kyverno-policies/) |
| **Golang** | [HTTP Server Kubernetes](./projects/golang/http-server-kubernetes/), [Kubernetes API Server IP](./projects/golang/kubernetes-api-server-ip/), [Unused Secret Detector](./projects/golang/unused-secret/) |
| **Infrastructure as Code** | [Terraform Module: Ingress NGINX](./projects/iac/terraform-modules/ingress-nginx/) |
| **Kubernetes** | [Amazon EKS Cluster](./projects/kubernetes/eks/), [k3d Cluster Setup](./projects/kubernetes/k3d/) |
| **Observability** | [Prometheus, Grafana, Loki Stack](./projects/observability/prometheus-grafana-loki/) |

---

## Core Focus Areas

- **Platform Reliability:** building and operating resilient Kubernetes platforms.
- **Delivery Automation:** implementing GitOps and CI/CD workflows for consistent releases.
- **Observability Engineering:** enabling metrics, logging, dashboards, and alerting for faster troubleshooting.
- **Security & Governance:** enforcing Kubernetes policy controls and improving secret management.

---

## Repository Structure

- [`projects/devops/`](./projects/devops/) — DevOps delivery and platform automation projects.
- [`projects/devsecops/`](./projects/devsecops/) — policy enforcement and security-focused Kubernetes projects.
- [`projects/golang/`](./projects/golang/) — Go-based applications and Kubernetes utilities.
- [`projects/iac/`](./projects/iac/) — infrastructure as code modules and reusable Terraform components.
- [`projects/kubernetes/`](./projects/kubernetes/) — Kubernetes platform and cluster setup projects.
- [`projects/observability/`](./projects/observability/) — monitoring, logging, and alerting stacks.
- [`.github/workflows/`](./.github/workflows/) — CI workflow examples.