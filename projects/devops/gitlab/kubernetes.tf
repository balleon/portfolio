resource "kubernetes_namespace_v1" "gitlab" {
  metadata {
    name = var.namespace
  }
}

# global.serviceAccount.create must be false in the Helm values: the chart
# refuses to let every subchart (webservice/sidekiq/gitaly/toolbox/...)
# create its own ServiceAccount under the same shared, custom name, so this
# single object is created here instead and shared across all of them.
resource "kubernetes_service_account_v1" "gitlab" {
  metadata {
    name      = "gitlab"
    namespace = kubernetes_namespace_v1.gitlab.metadata[0].name
    annotations = {
      "eks.amazonaws.com/role-arn" = aws_iam_role.gitlab.arn
    }
  }
}

resource "kubernetes_secret_v1" "db_password" {
  metadata {
    name      = "gitlab-postgresql-password"
    namespace = kubernetes_namespace_v1.gitlab.metadata[0].name
  }

  data = {
    password = random_password.db.result
  }
}

resource "kubernetes_secret_v1" "redis_password" {
  metadata {
    name      = "gitlab-redis-password"
    namespace = kubernetes_namespace_v1.gitlab.metadata[0].name
  }

  data = {
    password = random_password.redis.result
  }
}

# Consolidated object storage connection, shared by every S3-backed GitLab
# component (lfs/artifacts/uploads/packages/... and the toolbox backups).
# use_iam_profile lets each pod authenticate via the IRSA role instead of
# static access keys - see aws_iam_role.gitlab in s3.tf.
resource "kubernetes_secret_v1" "object_storage" {
  metadata {
    name      = "gitlab-rails-storage"
    namespace = kubernetes_namespace_v1.gitlab.metadata[0].name
  }

  data = {
    connection = yamlencode({
      provider        = "AWS"
      region          = data.aws_region.current.region
      use_iam_profile = true
    })
  }
}
