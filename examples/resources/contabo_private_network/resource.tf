# 1. Configure access to Fybe apis
provider "fybe" {
  oauth2_client_id     = "[your client id]"
  oauth2_client_secret = "[your client secret]"
  oauth2_user          = "[your username]"
  oauth2_pass          = "[your password]"
}

# 1. Create a virtual private cloud
resource "fybe_vpc" "k8s_vpc" {
  name        = "k8s_vpc"
	description = "virtual private cloud for k8s cluster"
	region 		   = "us-east-1"
  instance_ids = [666, 420]
}

# 3. Add more compute instances to this vpc
resource "fybe_vpc" "k8s_vpc" {
  instance_ids = [666, 420, 69]
}
