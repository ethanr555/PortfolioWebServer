
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
  instance_type = "t3.micro"
  region        = "us-west-1"
  source_ami_filter {
    filters = {
      name                = "ubuntu/images/*ubuntu-jammy-22.04-amd64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["099720109477"]
  }
  ssh_username = "ubuntu"
}

build {
  name = var.ami-name
  sources = [
    "source.amazon-ebs.ubuntu"
  ]

  provisioner "file" {
    source      = "../docker.tar"
    destination = "/tmp/docker.tar"
  }

  provisioner "file" {
    source      = "../dump.sql"
    destination = "/tmp/dump.sql"
  }

  provisioner "shell" {
    script = "bootstrap.sh"
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