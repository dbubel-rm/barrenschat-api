variable "aws_region" {
  description = "The region where AWS operations will take place"
  default     = "us-east-1"
}

variable "environment" {
    description = "Should be 'test', 'staging' or 'production'"
}

variable "access_key" {}
variable "secret_key" {}