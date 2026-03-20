output "backend_repository_url" {
  description = "Full ECR URL for backend image pushes."
  value       = aws_ecr_repository.backend.repository_url
}

output "frontend_repository_url" {
  description = "Full ECR URL for frontend image pushes."
  value       = aws_ecr_repository.frontend.repository_url
}
