terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 6.8"
    }
  }

  backend "s3" {
    bucket = var.backend_bucket
    key    = var.backend_key
    region = "us-west-1"
  }
}

resource "aws_vpc" "portfoliowebserver_VPC" {
}

resource "aws_subnet" "private_sb" {
  vpc_id = aws_vpc.portfoliowebserver_VPC.id
}

resource "aws_internet_gateway" "portfoliowebserver_gw" {
  vpc_id = aws_vpc.portfoliowebserver_VPC.id
}

data "aws_ami" "image" {
  most_recent = true

  owners     = ["self"]
  name_regex = "/^${var.ami-name}$/gm"
}

resource "aws_instance" "webserver" {
  ami           = data.aws_ami.image
  subnet_id     = aws_subnet.private_sb.id
  instance_type = "t3.micro"

  depends_on = [aws_internet_gateway.portfoliowebserver_gw]
}


resource "aws_cloudfront_vpc_origin" "portfoliowebserver_vpc_origin" {
  vpc_origin_endpoint_config {
    name                   = "portfoliowebserver-vpc-origin"
    arn                    = aws_instance.webserver.arn
    http_port              = 80
    https_port             = 443
    origin_protocol_policy = "http-only"
    origin_ssl_protocols {
      items    = ["TLSv1.2"]
      quantity = 1
    }
  }
}

locals {
  ec2_origin_id = "PortfolioWebserver-EC2-OriginID"
}

resource "aws_acm_certificate" "portfoliowebserver_cert" {
  domain_name               = var.domain
  subject_alternative_names = var.alternate-domains
  validation_method         = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "validation_record" {
  for_each = {
    for dvo in aws_acm_certificate.portfoliowebserver_cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  allow_overwrite = true
  name            = each.value.name
  records         = [each.value.record]
  ttl             = 60
  type            = each.value.type
  zone_id         = var.external-zone-id
}

resource "aws_acm_certificate_validation" "portfoliowebserver_validated_cert" {
  certificate_arn         = aws_acm_certificate.portfoliowebserver_cert.arn
  validation_record_fqdns = [for record in aws_route53_record.validation_record : record.fqdn]
}

resource "aws_cloudfront_distribution" "portfoliowebserver_cf" {
  origin {
    domain_name = aws_instance.webserver.private_dns
    origin_id   = local.ec2_origin_id
  }
  default_cache_behavior {
    allowed_methods        = ["GET"]
    cached_methods         = ["GET"]
    target_origin_id       = aws_instance.webserver.id
    viewer_protocol_policy = "https-only"
  }
  restrictions {
    geo_restriction {
      restriction_type = "whitelist"
      locations        = ["US"]
    }
  }
  viewer_certificate {
    acm_certificate_arn = aws_acm_certificate_validation.portfoliowebserver_validated_cert.certificate_arn
  }
  enabled = true
}

resource "aws_route53_record" "portfoliowebserver_dns" {
  type    = "A"
  name    = ""
  zone_id = var.external-zone-id

  alias {
    name                   = aws_cloudfront_distribution.portfoliowebserver_cf.domain_name
    zone_id                = aws_cloudfront_distribution.portfoliowebserver_cf.hosted_zone_id
    evaluate_target_health = false
  }
}
resource "aws_route53_record" "portfoliowebserver_dns_www" {
  type    = "A"
  name    = "www"
  zone_id = var.external-zone-id

  alias {
    name                   = aws_cloudfront_distribution.portfoliowebserver_cf.domain_name
    zone_id                = aws_cloudfront_distribution.portfoliowebserver_cf.hosted_zone_id
    evaluate_target_health = false
  }
}

