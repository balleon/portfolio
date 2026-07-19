terraform {
  required_version = ">= 1.12.0"

  backend "local" {}

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.55.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "3.2.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "3.2.1"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.9.0"
    }
  }
}
