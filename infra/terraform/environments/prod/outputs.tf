output "backend_ecr_url" {
  description = "Backend ECR repository URL."
  value       = module.ecr.backend_repository_url
}

output "frontend_ecr_url" {
  description = "Frontend ECR repository URL."
  value       = module.ecr.frontend_repository_url
}

output "managed_instance_id" {
  description = "EC2 instance ID when management mode is enabled."
  value       = module.compute_host.instance_id
}

output "managed_instance_public_ip" {
  description = "Elastic IP for the EC2 instance when Elastic IP management is enabled."
  value       = module.compute_host.elastic_ip
}

output "rds_endpoint" {
  description = "RDS connection endpoint (host:port) when management mode is enabled."
  value       = module.rds_instance.db_endpoint
}
