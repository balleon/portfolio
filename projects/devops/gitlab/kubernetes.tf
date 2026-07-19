resource "kubernetes_namespace_v1" "gitlab" {
  metadata {
    name = var.namespace
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
