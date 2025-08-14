
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
  type = list(string)
}

variable "backend_bucket" {
  type      = string
  sensitive = true
}

variable "backend_key" {
  type      = string
  sensitive = true
}