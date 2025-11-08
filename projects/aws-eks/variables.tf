variable "environment" {
  default = "test"
  type    = string
}

variable "cluster_name" {
  default = "test"
  type    = string
}

variable "cluster_cidr" {
  default = "192.168.0.0/16"
  type    = string
}

variable "cluster_private_subnets" {
  default = ["192.168.1.0/24", "192.168.2.0/24", "192.168.3.0/24"]
  type    = list(any)
}

variable "cluster_public_subnets" {
  default = ["192.168.4.0/24", "192.168.5.0/24", "192.168.6.0/24"]
  type    = list(any)
}