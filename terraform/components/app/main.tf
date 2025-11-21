terraform {
  cloud {
    organization = "tejasc"
    workspaces {
      name = "template-go"
    }
  }

  required_version = ">= 1.6.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

module "app" {
  source            = "../../modules/app"
  project_id        = var.project_id
  region            = var.region
  environment       = var.environment
  database_tier     = var.database_tier
  authorized_cidr   = var.authorized_cidr
  database_user     = var.database_user
  database_password = var.database_password
  database_name     = var.database_name
}
