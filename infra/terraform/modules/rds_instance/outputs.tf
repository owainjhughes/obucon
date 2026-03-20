output "db_endpoint" {
  description = "RDS instance connection endpoint (host:port)."
  value       = var.create_db ? aws_db_instance.main[0].endpoint : ""
}

output "db_port" {
  description = "RDS instance port."
  value       = var.create_db ? aws_db_instance.main[0].port : null
}
