provider "aws" {
  region = "us-west-1"
}
provider "aws" {
  region = "us-east-1"
  alias  = "virginia"
}