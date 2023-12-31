---
page_title: "fybe-terraform"
subcategory: ""
description: |-
  Manage your Fybe infrastructure with our terraform provider.

---

# Fybe Provider

A terraform provider for managing infrastructure offered by [Fybe](https://fybe.com) like compute instance, vpc or S3 object-storage.  

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api` (String) The api endpoint is https://api.fybe.com.
- `oauth2_client_id` (String) Your oauth2 client id can be found in the [Customer Control Panel](https://new.fybe.com/account/security) under the menu item account secret.
- `oauth2_client_secret` (String) Your oauth2 client secret can be found in the [Customer Control Panel](https://new.fybe.com/account/security) under the menu item account secret.
- `oauth2_pass` (String) API Password (this is a new password which you'll set or change in the [Customer Control Panel](https://new.fybe.com/account/security) under the menu item account secret.)
- `oauth2_token_url` (String) The oauth2 token url is https://airlock.fybe.com/auth/realms/fybe/protocol/openid-connect/token.
- `oauth2_user` (String) API User (your email address to login to the [Customer Control Panel](https://new.fybe.com/account/security) under the menu item account secret.
