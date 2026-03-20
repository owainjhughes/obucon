locals {
  common_tags = {
    Project     = "obucon"
    Environment = "prod"
    ManagedBy   = "terraform"
  }
}

module "ecr" {
  source = "../../modules/ecr_repositories"

  backend_repository_name  = var.backend_repository_name
  frontend_repository_name = var.frontend_repository_name
  tags                     = local.common_tags
}

module "rds_ingress" {
  source = "../../modules/rds_ingress_rule"

  create_rule           = var.manage_rds_ingress_rule
  app_security_group_id = var.app_security_group_id
  rds_security_group_id = var.rds_security_group_id
}

module "compute_host" {
  source = "../../modules/compute_host"

  create_instance   = var.manage_ec2_instance
  instance_name     = var.instance_name
  ami_id            = var.ami_id
  instance_type     = var.instance_type
  key_name          = var.key_name
  subnet_id         = var.subnet_id
  security_group_id = var.app_security_group_id
  tags              = local.common_tags
}
