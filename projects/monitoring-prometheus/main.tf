provider "helm" {
  kubernetes = {
    config_path = "~/.kube/config"
  }
}

resource "helm_release" "kube_prometheus_stack" {
  name             = "kube-prometheus-stack"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  version          = "75.7.0"
  namespace        = "prometheus"
  create_namespace = true

  set = [
    {
      name  = "alertmanager.ingress.enabled"
      value = "true"
    },
    {
      name  = "alertmanager.ingress.hosts[0]"
      value = var.ingress_hostname
    },
    {
      name  = "alertmanager.ingress.paths[0]"
      value = "/alertmanager"
    },
    {
      name  = "alertmanager.alertmanagerSpec.routePrefix"
      value = "/alertmanager"
    },
    {
      name  = "grafana.ingress.enabled"
      value = "true"
    },
    {
      name  = "grafana.ingress.hosts[0]"
      value = var.ingress_hostname
    },
    {
      name  = "grafana.ingress.paths"
      value = "/grafana"
    },
    {
      name  = "grafana.grafana\\.ini.server.domain"
      value = var.ingress_hostname
    },
    {
      name  = "grafana.grafana\\.ini.server.root_url"
      value = "%(protocol)s://%(domain)s/grafana"
    },
    {
      name  = "grafana.grafana\\.ini.server.serve_from_sub_path"
      value = "true"
    },
    {
      name  = "prometheus.ingress.enabled"
      value = "true"
    },
    {
      name  = "prometheus.ingress.hosts[0]"
      value = var.ingress_hostname
    },
    {
      name  = "prometheus.ingress.paths[0]"
      value = "/prometheus"
    },
    {
      name  = "prometheus.prometheusSpec.routePrefix"
      value = "/prometheus"
    }
  ]
}