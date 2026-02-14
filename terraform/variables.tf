variable "env" {
  description = "Environment (local, stage, prod)"
  type        = string
}

variable "project" {
  description = "Project name"
  type        = string
  default     = "snack"
}

variable "region" {
  description = "AWS Region"
  type        = string
  default     = "us-east-1"
}

variable "bucket_name" {
  description = "S3 Bucket name for frontend"
  type        = string
}

variable "persona_state_bucket" {
  description = "S3 bucket that stores snackPersona Terraform state"
  type        = string
}

variable "persona_state_key" {
  description = "State key for snackPersona Terraform state"
  type        = string
}

variable "persona_state_region" {
  description = "AWS region for snackPersona Terraform state bucket"
  type        = string
  default     = "us-east-1"
}

variable "persona_state_endpoint" {
  description = "Optional S3 endpoint for remote state (use LocalStack URL for local)"
  type        = string
  default     = ""
}
