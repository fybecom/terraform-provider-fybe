# 1. Configure access to Fybe apis
provider "fybe" {
  oauth2_client_id     = "[your client id]"
  oauth2_client_secret = "[your client secret]"
  oauth2_user          = "[your username]"
  oauth2_pass          = "[your password]"
}

# 2. Create password secret
resource "fybe_secret" "main_password" {
  name        = "my_secret"
	type        = "password"
	value 		  = "Test432!"
}

# 3. Update password secret
resource "fybe_secret" "main_password" {
	value 		  = "Test1234!"
}
