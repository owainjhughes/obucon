output "bucket_name" {
  description = "S3 bucket name for frontend static files."
  value       = length(aws_s3_bucket.frontend) > 0 ? aws_s3_bucket.frontend[0].id : ""
}

output "cloudfront_domain" {
  description = "CloudFront distribution domain name."
  value       = length(aws_cloudfront_distribution.frontend) > 0 ? aws_cloudfront_distribution.frontend[0].domain_name : ""
}

output "cloudfront_distribution_id" {
  description = "CloudFront distribution ID (used for cache invalidation in CI)."
  value       = length(aws_cloudfront_distribution.frontend) > 0 ? aws_cloudfront_distribution.frontend[0].id : ""
}
