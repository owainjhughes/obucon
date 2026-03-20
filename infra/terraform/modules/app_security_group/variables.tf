variable "create_sg" {
  description = "Set true after importing the app security group into Terraform state."
  type        = bool
  default     = false
}

variable "vpc_id" {
  description = "VPC ID the security group belongs to."
  type        = string
  default     = ""
}

variable "sg_name" {
  description = "Name of the app security group."
  type        = string
  default     = "launch-wizard-1"
}

variable "rds_security_group_id" {
  description = "RDS security group ID used for the PostgreSQL egress rule."
  type        = string
  default     = ""
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
  default     = {}
}
