terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  backend "s3" {}
}

provider "aws" {
  region = var.region
}

locals {
  persona_state_common = {
    bucket = var.persona_state_bucket
    key    = var.persona_state_key
    region = var.persona_state_region
  }

  persona_state_local_overrides = var.persona_state_endpoint == "" ? {} : {
    endpoint                    = var.persona_state_endpoint
    skip_credentials_validation = true
    skip_metadata_api_check     = true
    force_path_style            = true
  }
}

data "terraform_remote_state" "persona" {
  backend = "s3"
  config  = merge(local.persona_state_common, local.persona_state_local_overrides)
}

module "frontend" {
  source      = "../../terraform/service/frontend"
  bucket_name = var.bucket_name
}

module "backend" {
  source     = "../../terraform/service/backend"
  env        = var.env
  table_name = data.terraform_remote_state.persona.outputs.table_name
  table_arn  = data.terraform_remote_state.persona.outputs.table_arn
}

output "api_endpoint" {
  value = module.backend.api_endpoint
}

output "frontend_bucket" {
  value = module.frontend.bucket_name
}

output "cloudfront_domain" {
  value = module.frontend.cloudfront_domain
}
