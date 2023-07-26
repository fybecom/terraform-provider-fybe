terraform {
  required_providers {
    fybe = {
      source = "fybe/fybe"
      version = "__CURRENT_VERSION__"
    }
  }
}

# Set the following environment variables:
#
# FYBE_OAUTH2_CLIENT_ID
# FYBE_OAUTH2_CLIENT_SECRET
# FYBE_OAUTH2_USER
# FYBE_OAUTH2_PASS
#
# and you are good to go
provider "fybe" {}


# Create a default fybe VPS instance
resource "fybe_instance" "default_instance" {}

# Output our newly created instances
output "default_instance_output" {
  description = "Our first default instance"
  value = fybe_instance.default_instance
}
