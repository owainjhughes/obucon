# Adoption-safe RDS PostgreSQL instance.
#
# The password is placed in ignore_changes so Terraform never tries to rotate
# it after the initial import.  Manage credential rotation separately (e.g.
# via AWS Secrets Manager).

resource "aws_db_instance" "main" {
  count = var.create_db ? 1 : 0

  identifier        = var.identifier
  engine            = "postgres"
  engine_version    = "16"
  instance_class    = var.instance_class
  allocated_storage = var.allocated_storage
  storage_type      = "gp2"
  storage_encrypted = true

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  vpc_security_group_ids = var.vpc_security_group_ids
  db_subnet_group_name   = var.db_subnet_group_name
  parameter_group_name   = var.parameter_group_name

  multi_az            = false
  publicly_accessible = true
  skip_final_snapshot = true

  backup_retention_period = 1

  tags = var.tags

  lifecycle {
    ignore_changes = [
      password,
      engine_version,
      maintenance_window,
      backup_window,
      backup_retention_period,
      deletion_protection,
      final_snapshot_identifier,
      snapshot_identifier,
      apply_immediately,
      copy_tags_to_snapshot,
      max_allocated_storage,
      performance_insights_enabled,
    ]
  }
}
