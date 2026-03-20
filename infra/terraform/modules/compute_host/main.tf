resource "aws_instance" "app" {
  count = var.create_instance ? 1 : 0

  ami                    = var.ami_id
  instance_type          = var.instance_type
  key_name               = var.key_name
  subnet_id              = var.subnet_id
  vpc_security_group_ids = [var.security_group_id]

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
