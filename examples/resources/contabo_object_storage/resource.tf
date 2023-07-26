# 1. Configure access to Fybe apis
provider "fybe" {
  oauth2_client_id     = "[your client id]"
  oauth2_client_secret = "[your client secret]"
  oauth2_user          = "[your username]"
  oauth2_pass          = "[your password]"
}

# 2. Create a object-storage in us-central
resource "fybe_object_storage" "main_object_storage" {
  region                   = "us-central-1"
	total_purchased_space_tb = 1
}

# 3. Update main_object_storage, enable autoscaling
resource "fybe_object_storage" "main_object_storage" {
  auto_scaling {
    state         = "enabled"
    size_limit_tb = 10
  }
}
