variable "github_url" {
  type        = string
  description = "GitHub URL for where you want to configure runners."
}

variable "github_token" {
  type        = string
  description = "GitHub personal access token."
  sensitive   = true
}
