
variable "ami-name" {
  type     = string
  nullable = false
  default  = "PortfolioWebserver-ami"
}

variable "domain" {
  type     = string
  nullable = false
}

variable "external-zone-id" {
  type      = string
  sensitive = true
}

variable "alternate-domains" {
  type = string
}

variable "vpc-origin-name" {
  type = string
}