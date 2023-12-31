---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "fybe_object_storage_bucket Resource - terraform-provider-fybe"
subcategory: ""
description: |-
  Manage buckets on your Fybe object-storage. With this resource you are able to manage your buckets the same way your are able to manage them in your Fybe cockpit.
---

# fybe_object_storage_bucket (Resource)

Manage buckets on your Fybe object-storage. With this resource you are able to manage your buckets the same way your are able to manage them in your Fybe cockpit.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of your bucket, consider the naming restriction https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-s3-bucket-naming-requirements.html.
- `object_storage_id` (String) The Fybe objectStorageId on which the bucket should be created.

### Optional

- `public_sharing` (Boolean) Choose the access to your bucket. You can not share it at all or share it publicly.

### Read-Only

- `creation_date` (String) The creation date of the bucket.
- `id` (String) object-storage Id.
- `public_sharing_link` (String) If your bucket is publicly shared, you can access it with this link.
