variable "aws_region" {
  description = "AWS region for production resources."
  type        = string
  default     = "eu-west-2"
}

variable "backend_repository_name" {
  description = "Backend ECR repository name."
  type        = string
  default     = "obucon"
}

variable "frontend_repository_name" {
  description = "Frontend ECR repository name."
  type        = string
  default     = "obucon-frontend"
}

variable "app_security_group_id" {
  description = "Security group ID attached to the app EC2 instance."
  type        = string
  default     = ""
}

variable "rds_security_group_id" {
  description = "Security group ID attached to the RDS instance."
  type        = string
  default     = ""
}

variable "manage_rds_ingress_rule" {
  description = "Set true after confirming app and RDS security group IDs."
  type        = bool
  default     = false
}

variable "manage_ec2_instance" {
  description = "Set true after importing the EC2 instance into Terraform state."
  type        = bool
  default     = false
}

variable "instance_name" {
  description = "Name tag for the EC2 instance."
  type        = string
  default     = "obucon"
}

variable "ami_id" {
  description = "AMI ID for EC2 management mode."
  type        = string
  default     = "ami-xxxxxxxx"
}

variable "instance_type" {
  description = "EC2 instance type."
  type        = string
  default     = "t3.micro"
}

variable "key_name" {
  description = "EC2 key pair name."
  type        = string
  default     = ""
}

variable "subnet_id" {
  description = "Subnet ID for EC2 instance."
  type        = string
  default     = ""
}
