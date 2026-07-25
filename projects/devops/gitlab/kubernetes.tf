resource "kubernetes_namespace_v1" "gitlab" {
  metadata {
    name = var.namespace
  }
}

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
