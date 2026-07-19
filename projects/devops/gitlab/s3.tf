locals {
  s3_buckets = toset([
    "artifacts",
    "uploads",
    "lfs",
    "packages",
    "external-diffs",
    "terraform-state",
    "dependency-proxy",
    "backups",
    "registry",
  ])
}

resource "aws_s3_bucket" "this" {
  for_each = local.s3_buckets

  bucket = "${var.cluster_name}-gitlab-${each.key}-${data.aws_caller_identity.current.account_id}"

  tags = {
    Terraform = "true"
    Project   = "gitlab"
  }
}

resource "aws_s3_bucket_public_access_block" "this" {
  for_each = aws_s3_bucket.this

  bucket = each.value.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# IRSA role assumed by the GitLab ServiceAccount so pods reach S3 without
# static access keys. Trust is scoped to the ServiceAccount's namespace/name.
data "aws_iam_policy_document" "gitlab_irsa_trust" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [data.aws_iam_openid_connect_provider.this.arn]
    }

    condition {
      test     = "StringEquals"
      variable = "${replace(data.aws_iam_openid_connect_provider.this.url, "https://", "")}:sub"
      values   = ["system:serviceaccount:${var.namespace}:${var.service_account_name}"]
    }

    condition {
      test     = "StringEquals"
      variable = "${replace(data.aws_iam_openid_connect_provider.this.url, "https://", "")}:aud"
      values   = ["sts.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "gitlab" {
  name               = "${var.cluster_name}-gitlab-irsa"
  assume_role_policy = data.aws_iam_policy_document.gitlab_irsa_trust.json
}

data "aws_iam_policy_document" "gitlab_s3" {
  statement {
    effect    = "Allow"
    actions   = ["s3:ListBucket"]
    resources = [for b in aws_s3_bucket.this : b.arn]
  }

  statement {
    effect    = "Allow"
    actions   = ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"]
    resources = [for b in aws_s3_bucket.this : "${b.arn}/*"]
  }
}

resource "aws_iam_policy" "gitlab_s3" {
  name   = "${var.cluster_name}-gitlab-s3"
  policy = data.aws_iam_policy_document.gitlab_s3.json
}

resource "aws_iam_role_policy_attachment" "gitlab_s3" {
  role       = aws_iam_role.gitlab.name
  policy_arn = aws_iam_policy.gitlab_s3.arn
}
