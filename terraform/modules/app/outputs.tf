output "service_account_email" {
  value = local.service_account_email
}

output "database_instance" {
  value = module.postgres.instance_name
}

output "database_secret" {
  value = "${var.environment}-todo-database-url"
}

output "database_url" {
  value     = local.database_url
  sensitive = true
}
