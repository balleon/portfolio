resource "helm_release" "this" {
  name             = var.name
  namespace        = var.namespace
  repository       = "https://kubernetes.github.io/ingress-nginx"
  chart            = "ingress-nginx"
  version          = "4.14.0"
  create_namespace = true
  values           = var.values
}