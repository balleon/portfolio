resource "helm_release" "gitlab" {
  name       = "gitlab"
  repository = "https://charts.gitlab.io/"
  chart      = "gitlab"
  version    = "10.2.0"
  namespace  = kubernetes_namespace_v1.gitlab.metadata[0].name
  timeout    = 900

  # https://docs.gitlab.com/charts/charts/globals.html
  values = [
    yamlencode({
      # https://docs.gitlab.com/charts/charts/globals/
      global = {
        edition = "ce"

        # https://docs.gitlab.com/charts/charts/globals/#hostsuffix
        hosts = {
          domain = data.kubernetes_service_v1.traefik.status[0].load_balancer[0].ingress[0].hostname
          https  = false
          gitlab = {
            name  = data.kubernetes_service_v1.traefik.status[0].load_balancer[0].ingress[0].hostname
            https = false
          }
        }

        # https://docs.gitlab.com/charts/advanced/gateway-api/
        gatewayApi = {
          enabled              = false
          installEnvoy         = false
          configureCertmanager = false
        }

        # https://docs.gitlab.com/charts/charts/globals/#configure-ingress-settings
        ingress = {
          enabled              = true
          configureCertmanager = false
          class                = "traefik"
          tls = {
            enabled = false
          }
        }

        # https://docs.gitlab.com/charts/charts/globals/#service-accounts
        serviceAccount = {
          enabled = true
          create  = false
          name    = kubernetes_service_account_v1.gitlab.metadata[0].name
        }

        # https://docs.gitlab.com/charts/advanced/external-redis/#configure-the-chart
        redis = {
          host = "redis-master.${kubernetes_namespace_v1.gitlab.metadata[0].name}.svc.cluster.local"
          auth = {
            enabled = true
            secret  = kubernetes_secret_v1.redis_password.metadata[0].name
            key     = "password"
          }
        }

        # https://docs.gitlab.com/charts/advanced/external-db/
        psql = {
          host     = aws_db_instance.gitlab.address
          port     = aws_db_instance.gitlab.port
          username = var.db_username
          database = var.db_name
          password = {
            secret = kubernetes_secret_v1.db_password.metadata[0].name
            key    = "password"
          }
        }

        # https://docs.gitlab.com/charts/charts/globals/#configure-appconfig-settings
        appConfig = {
          object_store = {
            enabled = true
            connection = {
              secret = kubernetes_secret_v1.object_storage.metadata[0].name
              key    = "connection"
            }
          }
          lfs = {
            bucket = aws_s3_bucket.this["lfs"].id
          }
          artifacts = {
            bucket = aws_s3_bucket.this["artifacts"].id
          }
          uploads = {
            bucket = aws_s3_bucket.this["uploads"].id
          }
          packages = {
            bucket = aws_s3_bucket.this["packages"].id
          }
          externalDiffs = {
            enabled = true
            bucket  = aws_s3_bucket.this["external-diffs"].id
          }
          terraformState = {
            enabled = true
            bucket  = aws_s3_bucket.this["terraform-state"].id
          }
          dependencyProxy = {
            enabled = true
            bucket  = aws_s3_bucket.this["dependency-proxy"].id
          }
          backups = {
            bucket = aws_s3_bucket.this["backups"].id
          }
        }
      }

      installCertmanager = false

      "nginx-ingress" = {
        enabled = false
      }
      traefik = {
        install = false
      }
      prometheus = {
        install = false
      }
      "gitlab-runner" = {
        install = false
      }
      registry = {
        enabled = false
      }

      # Fixed single replica (min = max) for every component with an
      # HPA-driven or static replica count - demo traffic doesn't need
      # autoscaling. gitaly and gitlab-exporter have no such value: gitaly
      # replicates per named node rather than by scalar count, and
      # gitlab-exporter's replica count is hardcoded to 1 by the chart.
      gitlab = {
        webservice = {
          minReplicas = 1
          maxReplicas = 1
          resources = {
            requests = { cpu = "1", memory = "2Gi" }
            limits   = { cpu = "2", memory = "3Gi" }
          }
          workhorse = {
            resources = {
              requests = { cpu = "100m", memory = "100Mi" }
              limits   = { cpu = "300m", memory = "300Mi" }
            }
          }
        }

        sidekiq = {
          minReplicas = 1
          maxReplicas = 1
          resources = {
            requests = { cpu = "300m", memory = "1.5Gi" }
            limits   = { cpu = "900m", memory = "3Gi" }
          }
        }

        gitaly = {
          resources = {
            requests = { cpu = "100m", memory = "256Mi" }
            limits   = { cpu = "500m", memory = "1Gi" }
          }
        }

        "gitlab-shell" = {
          minReplicas = 1
          maxReplicas = 1
          resources = {
            requests = { cpu = "50m", memory = "64Mi" }
            limits   = { cpu = "200m", memory = "128Mi" }
          }
        }

        kas = {
          minReplicas = 1
          maxReplicas = 1
          resources = {
            requests = { cpu = "100m", memory = "256Mi" }
            limits   = { cpu = "300m", memory = "512Mi" }
          }
        }

        "gitlab-exporter" = {
          resources = {
            requests = { cpu = "75m", memory = "100Mi" }
            limits   = { cpu = "200m", memory = "200Mi" }
          }
        }

        migrations = {
          resources = {
            requests = { cpu = "500m", memory = "500Mi" }
            limits   = { cpu = "1", memory = "1Gi" }
          }
        }

        # https://docs.gitlab.com/charts/charts/gitlab/toolbox/
        toolbox = {
          replicas = 1
          resources = {
            requests = { cpu = "50m", memory = "350Mi" }
            limits   = { cpu = "500m", memory = "1Gi" }
          }
          backups = {
            objectStorage = {
              backend = "s3"
              config = {
                secret = kubernetes_secret_v1.object_storage.metadata[0].name
                key    = "connection"
              }
            }
          }
        }
      }
    })
  ]

  depends_on = [
    aws_db_instance.gitlab,
    helm_release.redis,
    kubernetes_service_account_v1.gitlab,
    kubernetes_secret_v1.db_password,
    kubernetes_secret_v1.object_storage,
    aws_iam_role_policy_attachment.gitlab_s3,
  ]
}
