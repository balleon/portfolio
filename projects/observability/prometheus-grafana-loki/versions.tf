terraform {
  required_version = ">= 1.12.0"

  backend "local" {}

  required_providers {
    helm = {
      source  = "hashicorp/helm"
      version = "3.0.2"
    }
  }
}