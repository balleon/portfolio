data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

# The EKS cluster is provisioned by another project (e.g. kubernetes/eks) and
# is only looked up here to source its VPC, subnets, security group and OIDC
# issuer for the resources below.
data "aws_eks_cluster" "this" {
  name = var.cluster_name
}

data "aws_iam_openid_connect_provider" "this" {
  url = data.aws_eks_cluster.this.identity[0].oidc[0].issuer
}
