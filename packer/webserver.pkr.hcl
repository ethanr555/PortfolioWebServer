
packer {
  required_plugins {
    amazon = {
      version = ">= 1.2.8"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

source "amazon-ebs" "ubuntu" {
  ami_name      = var.ami-name
  instance_type = "t4g.micro"
  region        = "us-west-1"
  source_ami_filter {
    filters = {
      name             = "ubuntu/images/*ubuntu-jammy-22.04-amd64-server-*"
      root-device-type = "ebs"
      virtualization   = "hvm"
    }
    most_recent = true
    owners      = ["${var.AWSAccountID}"]
  }
  ssh_username = "ubuntu"
}

build {
  name = var.ami-name
  sources = [
    "source.amazon-ebs.ubuntu"
  ]

  provisioner "shell" {
    script = bootstrap.sh
    environment_vars = [
      "SCRIPT_DBPASS=${var.DBPass}",
      "SCRIPT_DBUSER=${var.DBUser}",
      "SCRIPT_DBNAME=${var.DBName}",
      "SCRIPT_DBROOTPASS=${var.DBRootPass}",
      "SCRIPT_DUMPPATH=${var.DumpPath}",
      "SCRIPT_WEBSERVERDOCKERPATH=${var.WebServerDockerPath}"
    ]
  }
}