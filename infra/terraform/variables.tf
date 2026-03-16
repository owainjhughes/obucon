variable "aws_region" {
  description = "AWS region for all resources"
  type        = string
  default     = "eu-west-2"
}

variable "project_name" {
  description = "Project/application name used in tagging"
  type        = string
  default     = "obucon"
}

variable "vpc_id" {
  description = "Target VPC ID"
  type        = string
}

variable "subnet_id" {
  description = "Subnet ID where the EC2 instance is deployed"
  type        = string
}

variable "key_name" {
  description = "EC2 key pair name"
  type        = string
}

variable "ami_id" {
  description = "Optional custom AMI ID. Leave null to use latest Amazon Linux 2023"
  type        = string
  default     = null
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.micro"
}

variable "ssh_cidr" {
  description = "CIDR allowed to SSH into EC2"
  type        = string
}

variable "allowed_api_cidrs" {
  description = "CIDRs allowed to access backend API port"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "create_ec2_security_group" {
  description = "Create a new EC2 security group"
  type        = bool
  default     = true
}

variable "existing_ec2_security_group_id" {
  description = "Existing EC2 security group ID when create_ec2_security_group is false"
  type        = string
  default     = ""
}

variable "create_instance_profile" {
  description = "Create IAM role/profile with ECR read-only policy"
  type        = bool
  default     = true
}

variable "existing_instance_profile_name" {
  description = "Existing IAM instance profile name when create_instance_profile is false"
  type        = string
  default     = ""
}

variable "rds_security_group_id" {
  description = "RDS security group ID to authorize from EC2 security group"
  type        = string
}
