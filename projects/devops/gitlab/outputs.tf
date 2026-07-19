output "gitlab_url" {
  value       = "http://gitlab.${var.domain}"
  description = "URL of the GitLab web interface."
}

output "registry_url" {
  value       = "http://registry.${var.domain}"
  description = "URL of the GitLab container registry."
}

output "rds_endpoint" {
  value       = aws_db_instance.gitlab.address
  description = "RDS PostgreSQL endpoint used as global.psql.host."
}

output "elasticache_endpoint" {
  value       = aws_elasticache_cluster.gitlab.cache_nodes[0].address
  description = "ElastiCache Redis endpoint used as global.redis.host."
}

output "s3_buckets" {
  value       = { for k, b in aws_s3_bucket.this : k => b.id }
  description = "S3 buckets created for GitLab object storage."
}

output "irsa_role_arn" {
  value       = aws_iam_role.gitlab.arn
  description = "IAM role ARN assumed by the GitLab ServiceAccount for S3 access."
}
