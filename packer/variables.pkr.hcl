
variable "ami-name" {
  type    = string
  default = "PortfolioWebserver-ami"
}

variable "DBName" {
  type      = string
  sensitive = true
}

variable "DBUser" {
  type      = string
  sensitive = true
}

variable "DBPass" {
  type      = string
  sensitive = true
}

variable "DBRootPass" {
  type      = string
  sensitive = true
}

variable "DumpPath" {
  type = string
}

variable "WebServerDockerPath" {
  type = string
}

variable "AWSAccountID" {
  type      = string
  sensitive = true
}