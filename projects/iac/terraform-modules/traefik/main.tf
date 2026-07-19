resource "helm_release" "this" {
  name             = var.name
  namespace        = var.namespace
  repository       = "https://traefik.github.io/charts"
  chart            = "traefik"
  version          = "41.0.2"
  create_namespace = true
  values           = var.values
}