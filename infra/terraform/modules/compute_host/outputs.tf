output "instance_id" {
  description = "EC2 instance ID if managed by Terraform."
  value       = try(aws_instance.app[0].id, null)
}
