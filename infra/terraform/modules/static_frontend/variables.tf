variable "create" {
  description = "Set true to create the S3 bucket and CloudFront distribution."
  type        = bool
  default     = false
}

variable "bucket_name" {
  description = "Globally unique S3 bucket name for the frontend static files."
  type        = string
}

variable "custom_domain" {
  description = "Custom domain for CloudFront (e.g. obucon.com). Leave empty to use the default cloudfront.net domain."
  type        = string
  default     = ""
}

variable "tags" {
  description = "Tags to apply to resources."
  type        = map(string)
  default     = {}
}
