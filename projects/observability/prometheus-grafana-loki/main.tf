provider "helm" {
  kubernetes = {
    config_path = "~/.kube/config"
  }
}

resource "helm_release" "traefik" {
  name             = "traefik"
  repository       = "https://traefik.github.io/charts"
  chart            = "traefik"
  version          = "39.0.2"
  namespace        = "traefik"
  create_namespace = true

  set = [
    {
      name  = "metrics.prometheus.service.enabled"
      value = "true"
    }
  ]
}

resource "helm_release" "kube_prometheus_stack" {
  name             = "kube-prometheus-stack"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  version          = "82.4.3"
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
      name  = "prometheus.additionalServiceMonitors[0].endpoints[0].port"
      value = "metrics"
    },
    {
      name  = "prometheus.additionalServiceMonitors[0].name"
      value = "${helm_release.traefik.name}-monitor"
    },
    {
      name  = "prometheus.additionalServiceMonitors[0].namespaceSelector.matchNames[0]"
      value = helm_release.traefik.name
    },
    {
      name  = "prometheus.additionalServiceMonitors[0].selector.matchLabels.app\\.kubernetes\\.io/name"
      value = helm_release.traefik.name
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
    },
    {
      name  = "grafana.additionalDataSources[0].name"
      value = "Loki"
    },
    {
      name  = "grafana.additionalDataSources[0].access"
      value = "proxy"
    },
    {
      name  = "grafana.additionalDataSources[0].isDefault"
      value = "false"
    },
    {
      name  = "grafana.additionalDataSources[0].orgId"
      value = "1"
    },
    {
      name  = "grafana.additionalDataSources[0].type"
      value = "loki"
    },
    {
      name  = "grafana.additionalDataSources[0].url"
      value = "http://${helm_release.loki.name}-gateway.${helm_release.loki.namespace}.svc.cluster.local"
    }
  ]

  depends_on = [helm_release.traefik]
}

resource "helm_release" "loki" {
  name             = "loki"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "loki"
  version          = "6.53.0"
  namespace        = "loki"
  create_namespace = true

  set = [
    {
      name  = "loki.auth_enabled"
      value = "false"
    },
    {
      name  = "loki.commonConfig.replication_factor"
      value = "1"
    },
    {
      name  = "loki.schemaConfig.configs[0].from"
      value = "2024-04-01"
    },
    {
      name  = "loki.schemaConfig.configs[0].store"
      value = "tsdb"
    },
    {
      name  = "loki.schemaConfig.configs[0].object_store"
      value = "s3"
    },
    {
      name  = "loki.schemaConfig.configs[0].schema"
      value = "v13"
    },
    {
      name  = "loki.schemaConfig.configs[0].index.prefix"
      value = "loki_index_"
    },
    {
      name  = "loki.schemaConfig.configs[0].index.period"
      value = "24h"
    },
    {
      name  = "loki.pattern_ingester.enabled"
      value = "true"
    },
    {
      name  = "loki.limits_config.allow_structured_metadata"
      value = "true"
    },
    {
      name  = "loki.limits_config.volume_enabled"
      value = "true"
    },
    {
      name  = "loki.ruler.enable_api"
      value = "true"
    },
    {
      name  = "minio.enabled"
      value = "true"
    },
    {
      name  = "deploymentMode"
      value = "SingleBinary"
    },
    {
      name  = "singleBinary.replicas"
      value = "1"
    },
    {
      name  = "backend.replicas"
      value = "0"
    },
    {
      name  = "read.replicas"
      value = "0"
    },
    {
      name  = "write.replicas"
      value = "0"
    },
    {
      name  = "ingester.replicas"
      value = "0"
    },
    {
      name  = "querier.replicas"
      value = "0"
    },
    {
      name  = "queryFrontend.replicas"
      value = "0"
    },
    {
      name  = "queryScheduler.replicas"
      value = "0"
    },
    {
      name  = "distributor.replicas"
      value = "0"
    },
    {
      name  = "compactor.replicas"
      value = "0"
    },
    {
      name  = "indexGateway.replicas"
      value = "0"
    },
    {
      name  = "bloomCompactor.replicas"
      value = "0"
    },
    {
      name  = "bloomGateway.replicas"
      value = "0"
    },
    {
      name  = "chunksCache.enabled"
      value = "false"
    }
  ]
}

resource "helm_release" "promtail" {
  name             = "promtail"
  repository       = "https://grafana.github.io/helm-charts"
  chart            = "promtail"
  version          = "6.17.1"
  namespace        = "promtail"
  create_namespace = true

  set = [
    {
      name  = "config.clients[0].url"
      value = "http://${helm_release.loki.name}-gateway.${helm_release.loki.namespace}.svc.cluster.local/loki/api/v1/push"
    }
  ]
}