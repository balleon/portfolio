terraform {
  required_version = ">= 1.12.0"

  backend "s3" {}

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.100.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "3.1.0"
    }
  }
}