variable "cluster_name" {
  type        = string
  description = "Name of the existing EKS cluster to deploy GitLab into."
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
