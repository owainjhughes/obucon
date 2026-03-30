output "instance_id" {
  description = "EC2 instance ID if managed by Terraform."
  value       = try(aws_instance.app[0].id, null)
}

output "elastic_ip" {
  description = "Elastic IP address if managed by Terraform."
  value       = try(aws_eip.app[0].public_ip, null)
}

output "elastic_ip_allocation_id" {
  description = "Elastic IP allocation ID if managed by Terraform."
  value       = try(aws_eip.app[0].id, null)
}
