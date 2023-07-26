package fybe

import (
	"context"
	"strconv"

	apiClient "fybe.com/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func dataSourceSecret() *schema.Resource {
	return &schema.Resource{
		Description: "The Secret Management API offers the capability to store and handle passwords and SSH keys. Its usage is entirely optional and serves as a convenient feature, enabling easy reuse of SSH keys, among other functionalities.",
		ReadContext: dataSourceSecretRead,
		Schema: map[string]*schema.Schema{
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Secret creation date.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Secret modified date.",
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the secret.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Secret name.",
			},
			"value": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The secret's value will be accessible solely when retrieving an individual secret.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the secret. It will be available only when retrieving secrets, following types are allowed: `ssh`, `password`.",
			},
		},
	}
}

func dataSourceSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*apiClient.APIClient)

	var secretId int64
	var err error
	id := d.Get("id").(string)
	if id != "" {
		secretId, err = strconv.ParseInt(id, 10, 64)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	res, httpResp, err := client.SecretsApi.
		RetrieveSecret(ctx, secretId).
		XRequestId(uuid.NewV4().String()).
		Execute()

	if err != nil {
		return HandleResponseErrors(diags, httpResp)
	} else if len(res.Data) != 1 {
		return MultipleDataObjectsError(diags)
	}

	d.SetId(strconv.Itoa(int(res.Data[0].SecretId)))

	return AddSecretToData(res.Data[0], d, diags)
}
