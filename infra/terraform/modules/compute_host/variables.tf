variable "create_instance" {
  description = "Set true only when you are ready to manage/import the EC2 instance with Terraform."
  type        = bool
  default     = false
}

variable "instance_name" {
  description = "Name tag for the EC2 instance."
  type        = string
}

variable "ami_id" {
  description = "AMI ID used by the EC2 instance."
  type        = string
}

variable "instance_type" {
  description = "EC2 instance type."
  type        = string
  default     = "t3.micro"
}

variable "key_name" {
  description = "EC2 key pair name."
  type        = string
}

variable "subnet_id" {
  description = "Subnet ID for the EC2 instance."
  type        = string
}

variable "security_group_id" {
  description = "Security group ID attached to the EC2 instance."
  type        = string
}

variable "tags" {
  description = "Common tags."
  type        = map(string)
  default     = {}
}

variable "manage_elastic_ip" {
  description = "Set true to allocate and associate a dedicated Elastic IP to this instance."
  type        = bool
  default     = false
}

variable "desired_state" {
  description = "Set to running or stopped to control EC2 power state. Leave empty to skip state management."
  type        = string
  default     = ""

  validation {
    condition     = contains(["", "running", "stopped"], var.desired_state)
    error_message = "desired_state must be one of: \"\", \"running\", or \"stopped\"."
  }
}

variable "instance_profile_name" {
  description = "IAM instance profile name to attach to the EC2 instance. Leave empty to skip."
  type        = string
  default     = ""
}
