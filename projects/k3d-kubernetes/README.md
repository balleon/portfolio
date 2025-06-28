# Kubernetes cluster setup with NGINX Ingress Controller

This guide demonstrates how to create locally a lightweight Kubernetes cluster using [k3d](https://k3d.io/), install and configure the NGINX Ingress Controller. It also includes deploying a simple NGINX HTTP server and testing Ingress routes. It provisions:

- A **Kubernetes cluster** with **k3d**
- An **Ingress Controller** with **NGINX Ingress Controller**
- A **Kubernetes Deployment**, **Service**, and **Ingress** for **NGINX** HTTP server

## Prerequisites

- Access to a Linux or Unix-like terminal
- [Docker](https://docs.docker.com/engine/install/)
- [k3d](https://k3d.io/#installation)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Helm](https://helm.sh/docs/intro/install/)

## Usage

### 1. Create Kubernetes Cluster (Without Traefik)

By default, k3s comes with the Traefik Ingress Controller. The following command creates a cluster without it, and exposes HTTP (80) and HTTPS (443) via the k3d load balancer:

```bash
k3d cluster create test \
--port "80:80@loadbalancer" \
--port "443:443@loadbalancer" \
--k3s-arg "--disable=traefik@server:*"
```

### 2. Deploy NGINX Ingress Controller

Add the official Helm repository and install the NGINX Ingress Controller into a dedicated namespace:

```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx

helm install ingress-nginx ingress-nginx/ingress-nginx \
--namespace ingress-nginx \
--create-namespace \
--wait
```

### 3. Deploy NGINX HTTP Server with Ingress

Create a Namespace, deploy an NGINX pod, expose it internally, and define an Ingress rule:

```bash
kubectl create namespace nginx

kubectl create deployment nginx \
--namespace=nginx \
--image=nginx:1.29 \
--replicas=1

kubectl expose deployment nginx \
--namespace=nginx \
--port=8080 \
--target-port=80

kubectl create ingress nginx \
--namespace=nginx \
--class=nginx \
--rule="<server-private-ip>.nip.io/nginx*=nginx:8080"

kubectl annotate ingress nginx \
--namespace=nginx nginx.ingress.kubernetes.io/rewrite-target="/"
```

### 4. Test HTTP and HTTPS Routing

Use curl to test that NGINX Ingress Controller is properly routing traffic to NGINX HTTP server:

```bash
curl http://<server-private-ip>.nip.io/nginx
curl https://<server-private-ip>.nip.io/nginx --insecure
```

### 5. Cleanup

Tear down everything created in this guide:

```bash
kubectl delete ingress,service,deployment --namespace=nginx nginx
kubectl delete namespace nginx

helm uninstall ingress-nginx --namespace=ingress-nginx
helm repo remove ingress-nginx

k3d cluster delete test
```