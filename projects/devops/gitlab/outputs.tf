output "gitlab_url" {
  value       = "http://${data.kubernetes_service_v1.traefik.status[0].load_balancer[0].ingress[0].hostname}"
  description = "URL of the GitLab web interface."
}

output "rds_endpoint" {
  value       = aws_db_instance.gitlab.address
  description = "RDS PostgreSQL endpoint used as global.psql.host."
}

output "redis_endpoint" {
  value       = aws_elasticache_replication_group.gitlab.primary_endpoint_address
  description = "ElastiCache Redis endpoint used as global.redis.host."
}

output "s3_buckets" {
  value       = { for k, b in aws_s3_bucket.this : k => b.id }
  description = "S3 buckets created for GitLab object storage."
}
