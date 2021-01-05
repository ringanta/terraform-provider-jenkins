terraform {
  required_version = "= 0.13.5"

  required_providers {
    jenkins = {
      source  = "ringanta.id/ringanta/jenkins"
      version = "0.3"
    }
  }
}

provider "jenkins" {
  server_url = "http://localhost:8080"
  username   = "admin"
  password   = "adminpwd"
  ca_cert    = ""
}
