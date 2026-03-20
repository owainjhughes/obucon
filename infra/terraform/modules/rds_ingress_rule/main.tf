resource "aws_security_group_rule" "allow_postgres_from_app" {
  count = var.create_rule ? 1 : 0

  type                     = "ingress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  source_security_group_id = var.app_security_group_id
  security_group_id        = var.rds_security_group_id
  description              = "Allow app host to connect to PostgreSQL"
}
