# Kubernetes Cluster Setup with NGINX Ingress Controller

## Overview
This project creates a local Kubernetes cluster with k3d, installs NGINX Ingress Controller, and exposes an NGINX workload through Ingress.

## Goals
- Create a local k3d cluster with Traefik disabled.
- Install NGINX Ingress Controller using Helm.
- Deploy and expose an NGINX test application.

## Repository Structure
This project is command-driven and does not include manifest files; resources are created with `kubectl` and `helm` commands.

## Prerequisites
- Docker
- `k3d`
- `kubectl`
- `helm`

## Usage
### 1) Create cluster
```bash
k3d cluster create test \
--port "80:80@loadbalancer" \
--port "443:443@loadbalancer" \
--k3s-arg "--disable=traefik@server:*"
```

### 2) Install NGINX Ingress Controller
```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm install ingress-nginx ingress-nginx/ingress-nginx \
--namespace ingress-nginx \
--create-namespace \
--wait
```

### 3) Deploy NGINX and Ingress
```bash
kubectl create namespace nginx
kubectl create deployment nginx --namespace=nginx --image=nginx:1.29 --replicas=1
kubectl expose deployment nginx --namespace=nginx --port=8080 --target-port=80
kubectl create ingress nginx --namespace=nginx --class=nginx --rule="<server-private-ip>.nip.io/nginx*=nginx:8080"
kubectl annotate ingress nginx --namespace=nginx nginx.ingress.kubernetes.io/rewrite-target="/"
```

## Validation
```bash
curl http://<server-private-ip>.nip.io/nginx
curl https://<server-private-ip>.nip.io/nginx --insecure
```

## Cleanup
```bash
kubectl delete ingress,service,deployment --namespace=nginx nginx
kubectl delete namespace nginx
helm uninstall ingress-nginx --namespace=ingress-nginx
helm repo remove ingress-nginx
k3d cluster delete test
```