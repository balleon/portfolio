terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "3.0.1"
    }
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

resource "kubernetes_namespace_v1" "namespace" {
  for_each = toset(["ns-1", "ns-2"])

  metadata {
    name = each.key
  }
}

resource "kubernetes_service_account_v1" "service_account" {
  for_each = toset(["app1-sa", "app2-sa"])

  metadata {
    name      = each.key
    namespace = kubernetes_namespace_v1.namespace["ns-1"].metadata[0].name
  }
}