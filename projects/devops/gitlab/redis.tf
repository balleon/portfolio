resource "random_password" "redis" {
  length  = 24
  special = false
}

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

resource "aws_elasticache_replication_group" "gitlab" {
  replication_group_id = "${var.cluster_name}-gitlab"
  description          = "GitLab Redis"

  engine         = "redis"
  engine_version = "7.1"
  node_type      = "cache.t3.micro"

  num_cache_clusters = 1
  port               = 6379

  subnet_group_name  = aws_elasticache_subnet_group.gitlab.name
  security_group_ids = [aws_security_group.redis.id]

  transit_encryption_enabled = true
  at_rest_encryption_enabled = true
  auth_token                 = random_password.redis.result

  apply_immediately = true
}
