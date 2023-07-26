package fybe

import (
	"context"

	apiClient "fybe.com/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func dataSourceObjectStorage() *schema.Resource {
	return &schema.Resource{
		Description: "Create and manage our S3 compatible object-storage.  Please be aware that this resource is not the S3 API. If you wish to access the S3 API directly or use S3 compatible tools like the aws CLI, you can do so by using the S3 URL provided by this resource after creating. To retrieve the S3 credentials, please refer to the User Management API.",
		ReadContext: dataSourceObjectStorageRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Object Storage identifier",
			},
			"created_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "object-storage creation date.",
			},
			"cancel_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date on which the object-storage no longer available, due to cancellation.",
			},
			"s3_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "S3 URL needed to connect to the s3 API.",
			},
			"s3_tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The S3 tenant Id is only needed for public sharing.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The object-storage status. It can be set to `PROVISIONING|READY|UPGRADING|CANCELLED|ERROR|DISABLED`.",
			},
			"auto_scaling": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Autoscaling status it can be `enabled|disabled|error`.",
						},
						"size_limit_tb": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Optional:    true,
							Description: "Limit for the size to be autoscaled.",
						},
						"error_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "If the autoscaling is in an error state (see status property), the error message can be seen in this field.",
						},
					},
				},
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Your tenant Id.",
			},
			"customer_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Your customer Id.",
			},
			"data_center": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Data center located of this object-storage.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Region of this object-storage.",
			},
			"total_purchased_space_tb": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Amount of purchased space in terabyte.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name for this object-storage.",
			},
		},
	}
}

func dataSourceObjectStorageRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*apiClient.APIClient)

	var objectStorageId string
	var err error
	id := d.Get("id").(string)
	if id != "" {
		objectStorageId = id
	}

	if err != nil {
		return diag.FromErr(err)
	}

	res, httpResp, err := client.ObjectStoragesApi.RetrieveObjectStorage(ctx, objectStorageId).XRequestId(uuid.NewV4().String()).
		Execute()

	if err != nil {
		return HandleResponseErrors(diags, httpResp)
	} else if len(res.Data) != 1 {
		return MultipleDataObjectsError(diags)
	}

	d.SetId(res.Data[0].ObjectStorageId)

	return AddObjectStorageToData(
		res.Data[0],
		d,
		diags,
	)
}
