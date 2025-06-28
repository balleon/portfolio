variable "env" {
  type = string
}

variable "cluster_name" {
  type = string
}

variable "cluster_cidr" {
  type = string
}

variable "cluster_private_subnets" {
  type = list(any)
}

variable "cluster_public_subnets" {
  type = list(any)
}