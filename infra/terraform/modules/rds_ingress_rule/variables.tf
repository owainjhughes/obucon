variable "create_rule" {
  description = "Set true to manage the PostgreSQL ingress rule."
  type        = bool
  default     = true
}

variable "app_security_group_id" {
  description = "Security group attached to the application EC2 host."
  type        = string
  default     = ""
}

variable "rds_security_group_id" {
  description = "Security group attached to the RDS instance."
  type        = string
  default     = ""
}
