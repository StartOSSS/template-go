terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

locals {
  database_host         = coalesce(module.postgres.private_ip_address, module.postgres.public_ip_address)
  database_url          = "postgres://${var.database_user}:${var.database_password}@${local.database_host}:5432/${var.database_name}?sslmode=disable"
  service_account_id    = "${var.environment}-todo-api"
  service_account_email = "${local.service_account_id}@${var.project_id}.iam.gserviceaccount.com"
}

module "service_accounts" {
  source  = "terraform-google-modules/service-accounts/google"
  version = "~> 4.3"

  project_id   = var.project_id
  names        = [local.service_account_id]
  display_name = "${var.environment} todo api"
}

module "postgres" {
  source              = "terraform-google-modules/sql-db/google//modules/postgresql"
  version             = "~> 18.0"
  name                = "${var.environment}-todo"
  project_id          = var.project_id
  region              = var.region
  tier                = var.database_tier
  deletion_protection = false
  database_version    = "POSTGRES_15"
  availability_type   = "ZONAL"
  disk_autoresize     = true
  user_name     = var.database_user
  user_password = var.database_password
  db_name       = var.database_name
  ip_configuration = {
    ipv4_enabled    = true
    private_network = null
    authorized_networks = [{
      name  = "${var.environment}-cidr"
      value = var.authorized_cidr
    }]
  }
}

module "secrets" {
  source  = "GoogleCloudPlatform/secret-manager/google"
  version = "~> 0.4"

  project_id = var.project_id
  secrets = [{
    name        = "${var.environment}-todo-database-url"
    secret_data = local.database_url
  }]
}

resource "google_project_iam_member" "run_invoker" {
  project = var.project_id
  role    = "roles/run.invoker"
  member  = "serviceAccount:${local.service_account_email}"
}

resource "google_project_iam_member" "secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = "serviceAccount:${local.service_account_email}"
}
