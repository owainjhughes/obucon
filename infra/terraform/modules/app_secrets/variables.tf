variable "create" {
  description = "Set true to create the Secrets Manager secret and IAM instance profile."
  type        = bool
  default     = false
}

variable "secret_name" {
  description = "Name of the Secrets Manager secret (e.g. obucon/prod/app)."
  type        = string
  default     = "obucon/prod/app"
}

variable "role_name" {
  description = "Prefix for the IAM role and instance profile names."
  type        = string
  default     = "obucon"
}

variable "db_host" {
  description = "Database host."
  type        = string
  default     = ""
}

variable "db_port" {
  description = "Database port."
  type        = string
  default     = "5432"
}

variable "db_user" {
  description = "Database username."
  type        = string
  default     = "postgres"
}

variable "db_password" {
  description = "Database password."
  type        = string
  sensitive   = true
  default     = ""
}

variable "db_name" {
  description = "Database name."
  type        = string
  default     = "obucon"
}

variable "db_sslmode" {
  description = "PostgreSQL SSL mode."
  type        = string
  default     = "require"
}

variable "jwt_secret" {
  description = "JWT signing secret."
  type        = string
  sensitive   = true
  default     = ""
}

variable "allowed_origins" {
  description = "Comma-separated list of allowed CORS origins."
  type        = string
  default     = ""
}

variable "app_port" {
  description = "Port the backend listens on."
  type        = string
  default     = "8080"
}

variable "trusted_proxies" {
  description = "Comma-separated trusted proxy IPs."
  type        = string
  default     = "127.0.0.1,::1"
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
  default     = {}
}
