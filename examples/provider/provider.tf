terraform {
  required_providers {
    fybe = {
      source = "fybe/fybe"
      version = "1.0.0"
    }
  }
}

# Configure the authentication to Fybe, you can find every values you need in your Cockpit access and security settings -> https://cockpit.fybe.com/account/security
provider "fybe" {
  oauth2_user          = "[your username]"
  oauth2_pass          = "[your password]"
  oauth2_client_id     = "[your client id]"
  oauth2_client_secret = "[your client secret]"
}

# Create a first fybe VPS instance
resource "fybe_instance" "first_instance" {}

# Output our first created instance
output "first_instance_output" {
  description = "first instance"
  value       = fybe_instance.first_instance
}
