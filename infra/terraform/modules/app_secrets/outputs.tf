output "secret_arn" {
  description = "ARN of the Secrets Manager secret."
  value       = length(aws_secretsmanager_secret.app) > 0 ? aws_secretsmanager_secret.app[0].arn : ""
}

output "secret_name" {
  description = "Name of the Secrets Manager secret."
  value       = length(aws_secretsmanager_secret.app) > 0 ? aws_secretsmanager_secret.app[0].name : ""
}

output "instance_profile_name" {
  description = "IAM instance profile name to attach to the EC2 instance."
  value       = length(aws_iam_instance_profile.ec2) > 0 ? aws_iam_instance_profile.ec2[0].name : ""
}
