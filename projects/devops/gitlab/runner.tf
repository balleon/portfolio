resource "kubernetes_namespace_v1" "gitlab_runner" {
  metadata {
    name = "gitlab-runner"
  }
}

resource "kubernetes_service_account_v1" "gitlab_runner_token_provisioner" {
  metadata {
    name      = "gitlab-runner-token-provisioner"
    namespace = kubernetes_namespace_v1.gitlab.metadata[0].name
  }
}

resource "kubernetes_role_v1" "gitlab_runner_token_provisioner" {
  metadata {
    name      = "gitlab-runner-token-provisioner"
    namespace = kubernetes_namespace_v1.gitlab.metadata[0].name
  }

  rule {
    api_groups = [""]
    resources  = ["pods"]
    verbs      = ["get", "list"]
  }
  rule {
    api_groups = [""]
    resources  = ["pods/exec"]
    verbs      = ["create"]
  }
  rule {
    api_groups = [""]
    resources  = ["secrets"]
    verbs      = ["get", "create"]
  }
}

resource "kubernetes_role_binding_v1" "gitlab_runner_token_provisioner" {
  metadata {
    name      = "gitlab-runner-token-provisioner"
    namespace = kubernetes_namespace_v1.gitlab.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account_v1.gitlab_runner_token_provisioner.metadata[0].name
    namespace = kubernetes_service_account_v1.gitlab_runner_token_provisioner.metadata[0].namespace
  }

  role_ref {
    kind      = "Role"
    name      = kubernetes_role_v1.gitlab_runner_token_provisioner.metadata[0].name
    api_group = "rbac.authorization.k8s.io"
  }
}

resource "kubernetes_role_v1" "gitlab_runner_secret_writer" {
  metadata {
    name      = "gitlab-runner-secret-writer"
    namespace = kubernetes_namespace_v1.gitlab_runner.metadata[0].name
  }

  rule {
    api_groups = [""]
    resources  = ["secrets"]
    verbs      = ["get", "create"]
  }
}

resource "kubernetes_role_binding_v1" "gitlab_runner_secret_writer" {
  metadata {
    name      = "gitlab-runner-secret-writer"
    namespace = kubernetes_namespace_v1.gitlab_runner.metadata[0].name
  }

  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account_v1.gitlab_runner_token_provisioner.metadata[0].name
    namespace = kubernetes_service_account_v1.gitlab_runner_token_provisioner.metadata[0].namespace
  }

  role_ref {
    kind      = "Role"
    name      = kubernetes_role_v1.gitlab_runner_secret_writer.metadata[0].name
    api_group = "rbac.authorization.k8s.io"
  }
}

resource "kubernetes_job_v1" "gitlab_runner_token_provision" {
  metadata {
    name      = "gitlab-runner-token-provision"
    namespace = kubernetes_namespace_v1.gitlab.metadata[0].name
  }

  wait_for_completion = true

  timeouts {
    create = "10m"
    update = "10m"
  }

  spec {
    backoff_limit = 3

    template {
      metadata {}

      spec {
        service_account_name = kubernetes_service_account_v1.gitlab_runner_token_provisioner.metadata[0].name
        restart_policy       = "Never"

        container {
          name    = "provision"
          image   = "alpine/k8s:1.36.2"
          command = ["/bin/bash", "-c"]
          args = [
            <<-EOT
              set -euo pipefail

              RUNNER_NS="gitlab-runner"
              RUNNER_SECRET="gitlab-runner-secret"
              GITLAB_NS="gitlab"

              if kubectl get secret "$RUNNER_SECRET" --namespace="$RUNNER_NS" > /dev/null 2>&1; then
                echo "Secret $RUNNER_SECRET already exists in $RUNNER_NS, skipping."
                exit 0
              fi

              TOOLBOX_POD=$(kubectl get pod --namespace="$GITLAB_NS" \
                --selector=app=toolbox \
                --field-selector=status.phase=Running \
                --output=jsonpath='{.items[0].metadata.name}')

              if [ -z "$TOOLBOX_POD" ]; then
                echo "No running toolbox pod found in $GITLAB_NS Namespace."
                exit 1
              fi

              echo "Using toolbox Pod: $TOOLBOX_POD"

              TOKEN=$(kubectl exec --namespace="$GITLAB_NS" --container=toolbox "$TOOLBOX_POD" -- gitlab-rails runner "
                runner = Ci::Runner.find_or_create_by!(
                  runner_type: :instance_type,
                  description: 'k8s-autoscaling-runner'
                )
                print runner.token
              ")

              if [ -z "$TOKEN" ]; then
                echo "Failed to obtain runner token."
                exit 1
              fi

              kubectl create secret generic "$RUNNER_SECRET" \
                --namespace "$RUNNER_NS" \
                --from-literal=runner-token="$TOKEN" \
                --from-literal=runner-registration-token=""

              echo "Created secret $RUNNER_SECRET in $RUNNER_NS Namespace."
            EOT
          ]
        }
      }
    }
  }

  depends_on = [
    kubernetes_role_binding_v1.gitlab_runner_token_provisioner,
    kubernetes_role_binding_v1.gitlab_runner_secret_writer,
    helm_release.gitlab,
  ]
}

resource "helm_release" "gitlab_runner" {
  name       = "gitlab-runner"
  repository = "https://charts.gitlab.io/"
  chart      = "gitlab-runner"
  version    = "0.91.2"
  namespace  = kubernetes_namespace_v1.gitlab_runner.metadata[0].name

  values = [
    yamlencode({
      gitlabUrl = "http://${data.kubernetes_service_v1.traefik.status[0].load_balancer[0].ingress[0].hostname}"

      runners = {
        secret       = "gitlab-runner-secret"
        jobNamespace = kubernetes_namespace_v1.gitlab_runner.metadata[0].name
        config       = <<-EOT
          [[runners]]
            name = "gitlab-runner"
            executor = "kubernetes"
            [runners.kubernetes]
              namespace = "${kubernetes_namespace_v1.gitlab_runner.metadata[0].name}"
              image = "alpine:latest"
        EOT
      }

      rbac = {
        create = true
      }

      concurrent = 2

      resources = {
        requests = { cpu = "50m", memory = "64Mi" }
        limits   = { cpu = "200m", memory = "256Mi" }
      }
    })
  ]

  depends_on = [kubernetes_job_v1.gitlab_runner_token_provision]
}
