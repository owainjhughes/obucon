resource "aws_instance" "app" {
  count = var.create_instance ? 1 : 0

  ami                    = var.ami_id
  instance_type          = var.instance_type
  key_name               = var.key_name
  subnet_id              = var.subnet_id
  vpc_security_group_ids = [var.security_group_id]
  iam_instance_profile   = var.instance_profile_name != "" ? var.instance_profile_name : null

  tags = merge(var.tags, {
    Name = var.instance_name
  })

  # Adoption mode: avoid replacing an existing manually-created instance
  # while you progressively move settings into Terraform.
  lifecycle {
    ignore_changes = [
      ami,
      user_data,
      associate_public_ip_address,
      private_ip,
      root_block_device,
      ebs_block_device,
      metadata_options,
      credit_specification,
      monitoring,
      tags,
      volume_tags,
    ]
  }
}

resource "aws_ec2_instance_state" "app" {
  count = var.create_instance && var.desired_state != "" ? 1 : 0

  instance_id = aws_instance.app[0].id
  state       = var.desired_state
}

resource "aws_eip" "app" {
  count = var.create_instance && var.manage_elastic_ip ? 1 : 0

  domain = "vpc"

  tags = merge(var.tags, {
    Name = "${var.instance_name}-eip"
  })
}

resource "aws_eip_association" "app" {
  count = var.create_instance && var.manage_elastic_ip ? 1 : 0

  instance_id   = aws_instance.app[0].id
  allocation_id = aws_eip.app[0].id
}
