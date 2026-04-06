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

module "app_sg" {
  source = "../../modules/app_security_group"

  create_sg             = var.manage_app_sg
  vpc_id                = var.vpc_id
  sg_name               = var.app_sg_name
  rds_security_group_id = var.rds_security_group_id
  tags                  = local.common_tags
}

module "rds_ingress" {
  source = "../../modules/rds_ingress_rule"

  create_rule           = var.manage_rds_ingress_rule
  app_security_group_id = var.app_security_group_id
  rds_security_group_id = var.rds_security_group_id
}

module "rds_instance" {
  source = "../../modules/rds_instance"

  create_db              = var.manage_rds_instance
  identifier             = var.db_identifier
  instance_class         = var.db_instance_class
  allocated_storage      = var.db_allocated_storage
  db_name                = var.db_name
  db_username            = var.db_username
  db_password            = var.db_password
  vpc_security_group_ids = [var.rds_security_group_id]
  db_subnet_group_name   = var.db_subnet_group_name
  parameter_group_name   = var.db_parameter_group_name
  tags                   = local.common_tags
}

module "app_secrets" {
  source = "../../modules/app_secrets"

  create          = var.manage_app_secrets
  secret_name     = "obucon/prod/app"
  role_name       = "obucon-prod"
  db_host         = split(":", module.rds_instance.db_endpoint)[0]
  db_port         = "5432"
  db_user         = var.db_username
  db_password     = var.db_password
  db_name         = var.db_name
  db_sslmode      = "require"
  jwt_secret      = var.jwt_secret
  allowed_origins = var.allowed_origins
  trusted_proxies = "127.0.0.1,::1"
  tags            = local.common_tags
}

module "compute_host" {
  source = "../../modules/compute_host"

  create_instance      = var.manage_ec2_instance
  manage_elastic_ip    = var.manage_ec2_elastic_ip
  desired_state        = var.ec2_desired_state
  instance_name        = var.instance_name
  ami_id               = var.ami_id
  instance_type        = var.instance_type
  key_name             = var.key_name
  subnet_id            = var.subnet_id
  security_group_id    = var.app_security_group_id
  instance_profile_name = module.app_secrets.instance_profile_name
  tags                 = local.common_tags
}

module "static_frontend" {
  source = "../../modules/static_frontend"

  providers = {
    aws.us_east_1 = aws.us_east_1
  }

  create        = var.manage_static_frontend
  bucket_name   = var.frontend_bucket_name
  custom_domain = var.frontend_custom_domain
  tags          = local.common_tags
}
