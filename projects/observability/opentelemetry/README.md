# Deploy
k3d cluster create test --port "80:80@loadbalancer" --port "443:443@loadbalancer" --k3s-arg "--disable=traefik@server:*"
helmfile sync

kubectl apply -f manifests/collector.yaml
kubectl apply -f manifests/instrumentation.yaml

kubectl apply -f manifests/application.yaml

# Undeploy
kubectl delete -f manifests/
helmfile destroy
kubectl delete pvc --all --namespace=loki
k3d cluster delete test