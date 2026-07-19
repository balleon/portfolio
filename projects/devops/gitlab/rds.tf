resource "random_password" "db" {
  length  = 24
  special = false
}

resource "aws_db_subnet_group" "gitlab" {
  name       = "${var.cluster_name}-gitlab-db"
  subnet_ids = data.aws_eks_cluster.this.vpc_config[0].subnet_ids
}

resource "aws_security_group" "db" {
  name        = "${var.cluster_name}-gitlab-db"
  description = "Allow PostgreSQL access from the EKS cluster"
  vpc_id      = data.aws_eks_cluster.this.vpc_config[0].vpc_id

  ingress {
    description     = "PostgreSQL from the EKS cluster"
    from_port       = 5432
    to_port         = 5432
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

resource "aws_db_instance" "gitlab" {
  identifier = "${var.cluster_name}-gitlab"

  engine         = "postgres"
  engine_version = var.db_engine_version

  instance_class    = var.db_instance_class
  allocated_storage = var.db_allocated_storage
  storage_type      = "gp3"

  db_name  = var.db_name
  username = var.db_username
  password = random_password.db.result
  port     = 5432

  db_subnet_group_name   = aws_db_subnet_group.gitlab.name
  vpc_security_group_ids = [aws_security_group.db.id]

  # Single-AZ, no deletion protection: this is a demo/portfolio deployment,
  # not tuned for production availability guarantees.
  multi_az                = false
  publicly_accessible     = false
  backup_retention_period = 1
  apply_immediately       = true
  skip_final_snapshot     = true
  deletion_protection     = false
}
