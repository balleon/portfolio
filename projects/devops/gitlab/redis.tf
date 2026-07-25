resource "random_password" "redis" {
  length  = 24
  special = false
}

# The GitLab chart dropped its bundled Redis subchart in v10.0.0, so Redis is
# deployed here as a separate, minimal, standalone (single-node) release and
# wired into GitLab via global.redis in main.tf.
resource "helm_release" "redis" {
  name       = "redis"
  repository = "https://charts.bitnami.com/bitnami"
  chart      = "redis"
  version    = "27.0.17"
  namespace  = kubernetes_namespace_v1.gitlab.metadata[0].name

  values = [
    yamlencode({
      fullnameOverride = "redis"
      architecture     = "standalone"
      auth = {
        existingSecret            = kubernetes_secret_v1.redis_password.metadata[0].name
        existingSecretPasswordKey = "password"
      }
      master = {
        # Smallest built-in resource preset (100m/128Mi requests,
        # 150m/192Mi limits) - fine for demo traffic, not production.
        resourcesPreset = "nano"
      }
    })
  ]

  depends_on = [kubernetes_secret_v1.redis_password]
}
