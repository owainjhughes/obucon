variable "backend_repository_name" {
  description = "Name of the backend ECR repository."
  type        = string
}

variable "frontend_repository_name" {
  description = "Name of the frontend ECR repository."
  type        = string
}

variable "tags" {
  description = "Tags applied to both repositories."
  type        = map(string)
  default     = {}
}
