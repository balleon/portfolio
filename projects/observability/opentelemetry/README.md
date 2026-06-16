# Deploy
k3d cluster create test --port "80:80@loadbalancer" --port "443:443@loadbalancer" --k3s-arg "--disable=traefik@server:*"
helmfile sync

kubectl apply --filename=manifests/collector.yaml
kubectl apply --filename=manifests/instrumentation.yaml
kubectl apply --filename=manifests/deployment.yaml

curl http://172.22.142.254.nip.io/test

# URL
- http://172.22.142.254.nip.io/grafana
- http://172.22.142.254.nip.io/prometheus (`count by (__name__) ({job="test/test"})`)

# Undeploy
kubectl delete --filename=manifests/
helmfile destroy
kubectl delete pvc --all --namespace=loki
k3d cluster delete test