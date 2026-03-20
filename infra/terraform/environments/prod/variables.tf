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

variable "manage_app_sg" {
  description = "Set true after importing the app security group into Terraform state."
  type        = bool
  default     = false
}

variable "vpc_id" {
  description = "VPC ID for the app security group."
  type        = string
  default     = ""
}

variable "app_sg_name" {
  description = "Name of the app EC2 security group."
  type        = string
  default     = "launch-wizard-1"
}

variable "manage_rds_instance" {
  description = "Set true after importing the RDS instance into Terraform state."
  type        = bool
  default     = false
}

variable "db_identifier" {
  description = "RDS DB instance identifier."
  type        = string
  default     = "obucon-db"
}

variable "db_instance_class" {
  description = "RDS instance class."
  type        = string
  default     = "db.t3.micro"
}

variable "db_allocated_storage" {
  description = "Allocated storage in GiB for the RDS instance."
  type        = number
  default     = 20
}

variable "db_name" {
  description = "Name of the initial database."
  type        = string
  default     = "obucon"
}

variable "db_username" {
  description = "Master username for the RDS instance."
  type        = string
  default     = "postgres"
}

variable "db_password" {
  description = "Master password for the RDS instance (sensitive, gitignored via tfvars)."
  type        = string
  sensitive   = true
  default     = "placeholder-ignored-after-import"
}

variable "db_subnet_group_name" {
  description = "DB subnet group name."
  type        = string
  default     = ""
}

variable "db_parameter_group_name" {
  description = "DB parameter group name."
  type        = string
  default     = "default.postgres16"
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
