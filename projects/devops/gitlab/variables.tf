variable "cluster_name" {
  type        = string
  default     = "test"
  description = "Name of the existing EKS cluster to deploy GitLab into."
}

variable "namespace" {
  type        = string
  default     = "gitlab"
  description = "Kubernetes namespace the GitLab release is installed into."
}

variable "db_username" {
  type        = string
  default     = "gitlab"
  description = "Master username for the RDS PostgreSQL instance."
}

variable "db_name" {
  type        = string
  default     = "gitlab"
  description = "Database name created on the RDS PostgreSQL instance."
}
