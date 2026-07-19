resource "aws_elasticache_subnet_group" "gitlab" {
  name       = "${var.cluster_name}-gitlab-redis"
  subnet_ids = data.aws_eks_cluster.this.vpc_config[0].subnet_ids
}

resource "aws_security_group" "redis" {
  name        = "${var.cluster_name}-gitlab-redis"
  description = "Allow Redis access from the EKS cluster"
  vpc_id      = data.aws_eks_cluster.this.vpc_config[0].vpc_id

  ingress {
    description     = "Redis from the EKS cluster"
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [data.aws_eks_cluster.this.vpc_config[0].cluster_security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Single-node cache cluster, no AUTH token / in-transit encryption: this is a
# demo/portfolio deployment kept simple, matching the repo's HTTP-only demos.
resource "aws_elasticache_cluster" "gitlab" {
  cluster_id     = "${var.cluster_name}-gitlab"
  engine         = "redis"
  engine_version = var.redis_engine_version
  node_type      = var.redis_node_type

  num_cache_nodes = 1
  port            = 6379

  subnet_group_name  = aws_elasticache_subnet_group.gitlab.name
  security_group_ids = [aws_security_group.redis.id]

  apply_immediately = true
}
