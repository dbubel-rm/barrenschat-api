resource "aws_s3_bucket" "app" {
  bucket = "barrenschat"
  acl    = "private"
}

resource "aws_s3_bucket" "archive" {
  bucket = "archive12341234"
  acl    = "private"
}
