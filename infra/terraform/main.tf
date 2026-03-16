locals {
  common_tags = {
    Project   = var.project_name
    ManagedBy = "Terraform"
  }
}

data "aws_ami" "al2023" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-2023*-x86_64"]
  }
}

resource "aws_security_group" "ec2" {
  count       = var.create_ec2_security_group ? 1 : 0
  name        = "${var.project_name}-ec2-sg"
  description = "EC2 security group for ${var.project_name}"
  vpc_id      = var.vpc_id

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.ssh_cidr]
  }

  ingress {
    description = "Backend API"
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = var.allowed_api_cidrs
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-ec2-sg"
  })
}

resource "aws_iam_role" "ec2_role" {
  count = var.create_instance_profile ? 1 : 0
  name  = "${var.project_name}-ec2-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
  })

  tags = local.common_tags
}

resource "aws_iam_role_policy_attachment" "ecr_readonly" {
  count      = var.create_instance_profile ? 1 : 0
  role       = aws_iam_role.ec2_role[0].name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

resource "aws_iam_instance_profile" "ec2_profile" {
  count = var.create_instance_profile ? 1 : 0
  name  = "${var.project_name}-ec2-profile"
  role  = aws_iam_role.ec2_role[0].name
}

locals {
  ec2_sg_id = var.create_ec2_security_group ? aws_security_group.ec2[0].id : var.existing_ec2_security_group_id
  ec2_profile_name = var.create_instance_profile ? aws_iam_instance_profile.ec2_profile[0].name : var.existing_instance_profile_name
}

resource "aws_instance" "app" {
  ami                    = coalesce(var.ami_id, data.aws_ami.al2023.id)
  instance_type          = var.instance_type
  subnet_id              = var.subnet_id
  key_name               = var.key_name
  vpc_security_group_ids = [local.ec2_sg_id]
  iam_instance_profile   = local.ec2_profile_name

  user_data = <<-EOF
              #!/bin/bash
              set -euxo pipefail
              dnf update -y
              dnf install -y docker
              systemctl enable --now docker
              usermod -aG docker ec2-user
              EOF

  tags = merge(local.common_tags, {
    Name = "${var.project_name}-ec2"
  })
}

resource "aws_security_group_rule" "rds_ingress_from_ec2" {
  type                     = "ingress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  security_group_id        = var.rds_security_group_id
  source_security_group_id = local.ec2_sg_id
  description              = "Allow PostgreSQL from ${var.project_name} EC2"
}
