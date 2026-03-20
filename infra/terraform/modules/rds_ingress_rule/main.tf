resource "aws_vpc_security_group_ingress_rule" "allow_postgres_from_app" {
  count = var.create_rule ? 1 : 0

  security_group_id            = var.rds_security_group_id
  referenced_security_group_id = var.app_security_group_id
  ip_protocol                  = "tcp"
  from_port                    = 5432
  to_port                      = 5432
  description                  = "Allow app host to connect to PostgreSQL"
}
