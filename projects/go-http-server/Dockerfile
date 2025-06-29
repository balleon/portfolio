FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o kube-version

FROM alpine:3.22

WORKDIR /root/

COPY --from=builder /app/kube-version .

CMD ["./kube-version"]