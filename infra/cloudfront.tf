# variable "access_key" {}
# variable "secret_key" {}

# provider "aws" {
#   access_key = "${var.access_key}"
#   secret_key = "${var.secret_key}"
#   region     = "us-east-1"
# }

locals {
  s3_origin_id = "myS3Origin"
}

resource "aws_cloudfront_distribution" "s3_distribution" {
  origin {
    domain_name = "${aws_s3_bucket.app.bucket_regional_domain_name}"
    origin_id   = "${local.s3_origin_id}"

    # s3_origin_config {
    #   origin_access_identity = "origin-access-identity/cloudfront/ABCDEFG1234567"
    # }
  }

  enabled = true

  default_cache_behavior {
    allowed_methods        = ["GET"]
    cached_methods         = ["GET"]
    target_origin_id       = "${local.s3_origin_id}"
    viewer_protocol_policy = "redirect-to-https"

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    # viewer_protocol_policy = "allow-all"
    min_ttl     = 0
    default_ttl = 3600
    max_ttl     = 86400
  }

  tags {
    Environment = "test"
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  restrictions {
    geo_restriction {
      restriction_type = "whitelist"
      locations        = ["US", "CA", "GB", "DE"]
    }
  }

  #   is_ipv6_enabled     = true
  #   comment             = "Some comment"
  #   default_root_object = "index.html"

  #   logging_config {
  #     include_cookies = false
  #     bucket          = "mylogs.s3.amazonaws.com"
  #     prefix          = "myprefix"
  #   }

  #   aliases = ["mysite.example.com", "yoursite.example.com"]

  #   default_cache_behavior {
  #     allowed_methods  = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
  #     cached_methods   = ["GET", "HEAD"]
  #     target_origin_id = "${local.s3_origin_id}"

  #     forwarded_values {
  #       query_string = false

  #       cookies {
  #         forward = "none"
  #       }
  #     }

  #     viewer_protocol_policy = "allow-all"
  #     min_ttl                = 0
  #     default_ttl            = 3600
  #     max_ttl                = 86400
  #   }

  # Cache behavior with precedence 0
  #   ordered_cache_behavior {
  #     path_pattern     = "/content/immutable/*"
  #     allowed_methods  = ["GET", "HEAD", "OPTIONS"]
  #     cached_methods   = ["GET", "HEAD", "OPTIONS"]
  #     target_origin_id = "${local.s3_origin_id}"

  #     forwarded_values {
  #       query_string = false
  #       headers      = ["Origin"]

  #       cookies {
  #         forward = "none"
  #       }
  #     }

  #     min_ttl                = 0
  #     default_ttl            = 86400
  #     max_ttl                = 31536000
  #     compress               = true
  #     viewer_protocol_policy = "redirect-to-https"
  #   }

  # Cache behavior with precedence 1
  #   ordered_cache_behavior {
  #     path_pattern     = "/content/*"
  #     allowed_methods  = ["GET", "HEAD", "OPTIONS"]
  #     cached_methods   = ["GET", "HEAD"]
  #     target_origin_id = "${local.s3_origin_id}"

  #     forwarded_values {
  #       query_string = false

  #       cookies {
  #         forward = "none"
  #       }
  #     }

  #     min_ttl                = 0
  #     default_ttl            = 3600
  #     max_ttl                = 86400
  #     compress               = true
  #     viewer_protocol_policy = "redirect-to-https"
  #   }

  #   price_class = "PriceClass_200"

  #   restrictions {
  #     geo_restriction {
  #       restriction_type = "whitelist"
  #       locations        = ["US", "CA", "GB", "DE"]
  #     }
  #   }
}

# resource "aws_s3_bucket" "testbucket" {
#   bucket = "cohesionapps"
#   acl    = "private"


# #   tags {
# #     Name = "My bucket"
# #   }
# }


# resource "aws_instance" "example" {
#   ami           = "ami-b374d5a5"
#   instance_type = "t2.micro"
# }
# resource "aws_eip" "ip" {
#   instance = "${aws_instance.example.id}"
# }

