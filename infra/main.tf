# IMPORTANT: Make sure the value of backend.key below matches the environment (test, staging or production)

#   backend "s3" {
#     bucket = "cohesion-client-js-backend"
#     # key    = "test/terraform.tfstate"
#     region = "us-east-1"
#   }

provider "aws" {
  # version = "~> 1.35"
  # shared_credentials_file = "$HOME/.aws/credentials"
  profile = "default"

  # region = "${backend.region}"
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "us-east-1"
}

create-bucket --region us-east-1 --create-bucket-configuration LocationConstraint=us-east-1 --acl private --bucket cohesion-clients-terraform-artifacts