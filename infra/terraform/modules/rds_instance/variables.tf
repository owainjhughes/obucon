variable "create_db" {
  description = "Set true after importing the RDS instance into Terraform state."
  type        = bool
  default     = false
}

variable "identifier" {
  description = "RDS DB instance identifier."
  type        = string
  default     = "obucon-db"
}

variable "instance_class" {
  description = "RDS instance class."
  type        = string
  default     = "db.t3.micro"
}

variable "allocated_storage" {
  description = "Allocated storage in GiB."
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
  description = "Master password for the RDS instance (managed outside Terraform after import)."
  type        = string
  sensitive   = true
  default     = "placeholder-ignored-after-import"
}

variable "vpc_security_group_ids" {
  description = "List of VPC security group IDs to attach to the RDS instance."
  type        = list(string)
  default     = []
}

variable "db_subnet_group_name" {
  description = "DB subnet group name."
  type        = string
  default     = ""
}

variable "parameter_group_name" {
  description = "DB parameter group name."
  type        = string
  default     = "default.postgres16"
}

variable "tags" {
  description = "Tags to apply to the RDS instance."
  type        = map(string)
  default     = {}
}
