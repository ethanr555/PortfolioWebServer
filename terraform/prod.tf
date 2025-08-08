
variable "domain_name" {}

resource "digitalocean_droplet" "portfoliowebserver-webserver" {
  image = "ubuntu-20-04-x64"
  name = "portfoliowebserver-webserver"
  region = "sfo3"
  size = "s-1vcpu-512mb-10gb"
}

resource "digitalocean_droplet" "portfoliowebserver-database" {
    image = "ubuntu-20-04-x64"
    name = "portfoliowebserver-database"
    region = "sfo3"
    size = "s-1vcpu-512mb-10gb"
}

resource "digitalocean_domain" "domain" {
  name = var.domain_name
}

resource "digitalocean_record" "portfoliowebserver-record-a" {
    domain = digitalocean_domain.domain
    type = "A"
    name = ""
    value = digitalocean_droplet.portfoliowebserver-webserver.ipv4_address
}

resource "digitalocean_record" "portfoliowebserver-record-cname" {
    domain = digitalocean_domain.domain
    type = "CNAME"
    name = "www"
    value = digitalocean_record.portfoliowebserver-record-a
}

resource "digitalocean_record" "ns1" {
  domain = digitalocean_domain.domain
  type = "NS"
  name = ""
  value = "ns1.digitalocean.com"
}

resource "digitalocean_record" "ns2" {
  domain = digitalocean_domain.domain
  type = "NS"
  name = ""
  value = "ns2.digitalocean.com"
}

resource "digitalocean_record" "ns3" {
  domain = digitalocean_domain.domain
  type = "NS"
  name = ""
  value = "ns3.digitalocean.com"
}
