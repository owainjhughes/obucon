resource "aws_secretsmanager_secret" "app" {
  count = var.create ? 1 : 0

  name                    = var.secret_name
  description             = "Application secrets for ${var.secret_name}"
  recovery_window_in_days = 0

  tags = var.tags
}

resource "aws_secretsmanager_secret_version" "app" {
  count = var.create ? 1 : 0

  secret_id = aws_secretsmanager_secret.app[0].id

  secret_string = jsonencode({
    DB_HOST         = var.db_host
    DB_PORT         = var.db_port
    DB_USER         = var.db_user
    DB_PASSWORD     = var.db_password
    DB_NAME         = var.db_name
    DB_SSLMODE      = var.db_sslmode
    JWT_SECRET      = var.jwt_secret
    ALLOWED_ORIGINS = var.allowed_origins
    PORT            = var.app_port
    APP_ENV         = "production"
    COOKIE_SECURE   = "true"
    TRUSTED_PROXIES = var.trusted_proxies
  })

  lifecycle {
    ignore_changes = [secret_string]
  }
}

# IAM role for the EC2 instance to read secrets
resource "aws_iam_role" "ec2" {
  count = var.create ? 1 : 0

  name = "${var.role_name}-ec2-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })

  tags = var.tags
}

resource "aws_iam_role_policy" "secrets_read" {
  count = var.create ? 1 : 0

  name = "read-app-secrets"
  role = aws_iam_role.ec2[0].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = ["secretsmanager:GetSecretValue"]
        Resource = aws_secretsmanager_secret.app[0].arn
      },
      {
        Effect   = "Allow"
        Action   = ["ecr:GetAuthorizationToken",
                    "ecr:BatchCheckLayerAvailability",
                    "ecr:GetDownloadUrlForLayer",
                    "ecr:BatchGetImage"]
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_instance_profile" "ec2" {
  count = var.create ? 1 : 0

  name = "${var.role_name}-ec2-profile"
  role = aws_iam_role.ec2[0].name

  tags = var.tags
}
