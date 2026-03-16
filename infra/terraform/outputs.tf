output "ec2_instance_id" {
  description = "EC2 instance ID"
  value       = aws_instance.app.id
}

output "ec2_public_ip" {
  description = "EC2 public IP"
  value       = aws_instance.app.public_ip
}

output "ec2_security_group_id" {
  description = "EC2 security group ID in use"
  value       = local.ec2_sg_id
}

output "instance_profile_name" {
  description = "IAM instance profile attached to EC2"
  value       = local.ec2_profile_name
}
