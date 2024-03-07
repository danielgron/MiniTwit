terraform {
  required_version = ">= 1.0.0"

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

variable "do_token" {
  description = "DigitalOcean API Token"
  type        = string
}

variable "pvt_key" {}

provider "digitalocean" {
  token = var.do_token
}

resource "digitalocean_droplet" "prod" {
  image  = "docker-20-04"
  name   = "prod"
  region = "ams3"
  size   = "s-1vcpu-1gb"
  ssh_keys = [
    data.digitalocean_ssh_key.Viktoria_key.id
  ]

  connection {
    host = self.ipv4_address
    user = "root"
    type = "ssh"
    private_key = var.pvt_key
    timeout     = "2m"
  }

  provisioner "file" {
    source      = "docker_compose.yml"
    destination = "/tmp/docker_compose.yml"
  }

  provisioner "file" {
    source      = "provision.sh"
    destination = "/tmp/provision.sh"
  }

  provisioner "remote-exec" {
  inline = [
    "chmod +x /tmp/provision.sh", 
    "/tmp/provision.sh"            
  ]
  }
}

data "digitalocean_ssh_key" "Viktoria_key" {
  name = "Viktoria_key"
}

