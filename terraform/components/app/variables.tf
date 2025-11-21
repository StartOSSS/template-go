variable "project_id" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
}

variable "environment" {
  description = "Deployment environment label"
  type        = string
}

variable "database_tier" {
  description = "Cloud SQL tier"
  type        = string
}

variable "authorized_cidr" {
  description = "CIDR allowed to reach the database"
  type        = string
}

variable "database_user" {
  description = "Database user"
  type        = string
}

variable "database_name" {
  description = "Database name"
  type        = string
}
