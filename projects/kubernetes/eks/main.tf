module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.21.0"

  name = var.cluster_name
  cidr = var.cluster_cidr

  azs             = slice(data.aws_availability_zones.available.names, 0, 3)
  private_subnets = var.cluster_private_subnets
  public_subnets  = var.cluster_public_subnets

  enable_nat_gateway = true
  single_nat_gateway = true

  # https://docs.aws.amazon.com/eks/latest/userguide/network-load-balancing.html
  public_subnet_tags = {
    "kubernetes.io/role/elb" = 1
  }

  # https://docs.aws.amazon.com/eks/latest/userguide/network-load-balancing.html
  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = 1
  }

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "20.37.1"

  cluster_name                   = var.cluster_name
  cluster_version                = "1.33"
  cluster_endpoint_public_access = true

  enable_cluster_creator_admin_permissions = true

  cluster_compute_config = {
    enabled    = true
    node_pools = ["general-purpose", "system"]
  }

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  tags = {
    Terraform   = "true"
    Environment = var.environment
  }
}

module "ingress_nginx" {
  source = "../../iac/terraform-modules/ingress-nginx/"

  values = [
    <<-EOT
    controller:
      resources:
        limits:
          cpu: 100m
          memory: 100Mi
        requests:
          cpu: 100m
          memory: 100Mi
      service:
        annotations:
          service.beta.kubernetes.io/aws-load-balancer-scheme: "internet-facing"
    EOT
  ]
}