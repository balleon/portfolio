```
kubectl apply -f ./kubernetes.yaml

go mod init github.com/balleon/unused-secret
go mod tidy
go run main.go
```