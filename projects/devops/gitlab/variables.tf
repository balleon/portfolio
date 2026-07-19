variable "cluster_name" {
  type        = string
  description = "Name of the existing EKS cluster to deploy GitLab into."
}

variable "domain" {
  type        = string
  description = "Base domain used to derive GitLab hostnames (gitlab.<domain>, registry.<domain>)."
}

variable "namespace" {
  type        = string
  default     = "gitlab"
  description = "Kubernetes namespace the GitLab release is installed into."
}

variable "service_account_name" {
  type        = string
  default     = "gitlab"
  description = "Shared Kubernetes ServiceAccount name used by GitLab components to assume the S3 IRSA role."
}

variable "db_username" {
  type        = string
  default     = "gitlab"
  description = "Master username for the RDS PostgreSQL instance."
}

variable "db_name" {
  type        = string
  default     = "gitlabhq_production"
  description = "Database name created on the RDS PostgreSQL instance."
}

variable "db_engine_version" {
  type        = string
  default     = "16.4"
  description = "PostgreSQL engine version for RDS."
}

variable "db_instance_class" {
  type        = string
  default     = "db.t3.medium"
  description = "RDS instance class."
}

variable "db_allocated_storage" {
  type        = number
  default     = 20
  description = "Allocated storage (GiB) for the RDS instance."
}

variable "redis_node_type" {
  type        = string
  default     = "cache.t3.micro"
  description = "ElastiCache node type."
}

variable "redis_engine_version" {
  type        = string
  default     = "7.1"
  description = "Redis OSS engine version for ElastiCache."
}
