variable "name" {
  type        = string
  default     = "ingress-nginx"
  description = "Release name"
}

variable "namespace" {
  type        = string
  default     = "ingress-nginx"
  description = "Namespace to install the release into"
}

variable "values" {
  type        = list(string)
  default     = []
  description = "List of values in raw yaml format to pass to helm"
}