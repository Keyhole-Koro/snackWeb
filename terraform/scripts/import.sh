#!/usr/bin/env bash
set -euo pipefail

# Usage: ./scripts/import.sh <tfvars_file>
# Optional env vars:
#   LAMBDA_FUNCTION_NAME        (default: snack_backend-<env>)
#   CLOUDFRONT_DISTRIBUTION_ID  (if set, CloudFront is imported)
TFVARS_FILE=${1:-}
if [ -z "$TFVARS_FILE" ]; then
  echo "Usage: $0 <tfvars_file>"
  exit 1
fi

BUCKET_NAME=$(grep -E '^bucket_name' "$TFVARS_FILE" | cut -d'=' -f2 | tr -d ' "')
ENV=$(grep -E '^env' "$TFVARS_FILE" | cut -d'=' -f2 | tr -d ' "')
LAMBDA_FUNCTION_NAME=${LAMBDA_FUNCTION_NAME:-snack_backend-${ENV}}

if [ "$ENV" = "local" ]; then
  export AWS_ENDPOINT_URL="http://localhost:4566"
  export AWS_REGION="us-east-1"
  export AWS_ACCESS_KEY_ID="test"
  export AWS_SECRET_ACCESS_KEY="test"
fi

import_if_missing() {
  local addr=$1
  local id=$2

  if terraform state list | grep -Fq "$addr"; then
    echo "skip: $addr already in state"
    return 0
  fi

  echo "import: $addr <- $id"
  terraform import -var-file="$TFVARS_FILE" "$addr" "$id" || true
}

import_if_missing "module.frontend.aws_s3_bucket.frontend_bucket" "$BUCKET_NAME"
import_if_missing "module.backend.aws_lambda_function.api_backend" "$LAMBDA_FUNCTION_NAME"
import_if_missing "module.backend.aws_lambda_function_url.api_url" "$LAMBDA_FUNCTION_NAME"

if [ -n "${CLOUDFRONT_DISTRIBUTION_ID:-}" ]; then
  import_if_missing "module.frontend.aws_cloudfront_distribution.cdn" "$CLOUDFRONT_DISTRIBUTION_ID"
fi

echo "done"
