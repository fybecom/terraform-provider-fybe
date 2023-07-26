# 1. Configure access to Fybe apis
provider "fybe" {
  oauth2_client_id     = "[your client id]"
  oauth2_client_secret = "[your client secret]"
  oauth2_user          = "[your username]"
  oauth2_pass          = "[your password]"
}

# 2. Create a new compute instance with specs of the V12 product
resource "fybe_instance" "control_plane_instance" {
  display_name  = "control-plane"
  product_id    = "V12"
}

# 3. Update custom image on instance
resource "fybe_instance" "control_plane_instance" {
  image_id = fybe_image.custom_image_templeos.id
}
