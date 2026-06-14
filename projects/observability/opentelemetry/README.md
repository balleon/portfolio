k3d cluster create test --port "80:80@loadbalancer" --port "443:443@loadbalancer" --k3s-arg "--disable=traefik@server:*"
helmfile sync

helmfile destroy
kubectl delete pvc --all -n loki
k3d cluster delete test