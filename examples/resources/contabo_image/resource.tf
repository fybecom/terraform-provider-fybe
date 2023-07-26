# 1. Configure access to Fybe apis
provider "fybe" {
  oauth2_user          = "[your username]"
  oauth2_pass          = "[your password]"
  oauth2_client_id     = "[your client id]"
  oauth2_client_secret = "[your client secret]"
}

# 2. Create templeos custom image 
resource "fybe_image" "custom_image_templeos" {
  name        = "templeos"
  image_url   = "https://templeos.org/Downloads/TempleOSLite.ISO"
  os_type     = "Linux"
  version     = "lastest"
  description = "An idiot admires complexity, a genius admires simplicity"
}
