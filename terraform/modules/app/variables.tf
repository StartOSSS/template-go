variable "project_id" {
  type = string
}

variable "region" {
  type = string
}

variable "environment" {
  type = string
}

variable "database_tier" {
  type    = string
  default = "db-custom-1-3840"
}

variable "authorized_cidr" {
  type    = string
  default = "0.0.0.0/0"
}

variable "database_user" {
  type    = string
  default = "todo"
}

variable "database_name" {
  type    = string
  default = "todo"
}
