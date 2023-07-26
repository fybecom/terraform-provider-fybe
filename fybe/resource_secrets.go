package fybe

import (
	"context"
	"strconv"
	"time"

	openapi "fybe.com/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		Description:   "The Secret Management API offers the capability to store and handle passwords and SSH keys. Its usage is entirely optional and serves as a convenient feature, enabling easy reuse of SSH keys, among other functionalities.",
		CreateContext: resourceSecretCreate,
		ReadContext:   resourceSecretRead,
		UpdateContext: resourceSecretUpdate,
		DeleteContext: resourceSecretDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

func resourceSecretCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*openapi.APIClient)

	secretName := d.Get("name").(string)
	secretValue := d.Get("value").(string)
	secretType := d.Get("type").(string)

	createSecretRequest := openapi.NewCreateSecretRequestWithDefaults()
	createSecretRequest.Name = secretName
	createSecretRequest.Value = secretValue
	createSecretRequest.Type = secretType

	res, httpResp, err := client.SecretsApi.
		CreateSecret(context.Background()).
		XRequestId(uuid.NewV4().String()).
		CreateSecretRequest(*createSecretRequest).
		Execute()

	if err != nil {
		return HandleResponseErrors(diags, httpResp)
	}

	if len(res.Data) != 1 {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Internal Error: should have returned only one object",
		})
	}

	d.SetId(strconv.Itoa(int(res.Data[0].SecretId)))
	return resourceSecretRead(ctx, d, m)
}

func resourceSecretRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*openapi.APIClient)

	secretId, err := strconv.ParseInt(d.Id(), 10, 64)

	if err != nil {
		return diag.FromErr(err)
	}

	res, httpResp, err := client.SecretsApi.
		RetrieveSecret(ctx, secretId).
		XRequestId(uuid.NewV4().String()).
		Execute()

	if err != nil {
		return HandleResponseErrors(diags, httpResp)
	}

	if len(res.Data) != 1 {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Internal Error: should have returned only one object",
		})
	}

	return AddSecretToData(res.Data[0], d, diags)
}

func resourceSecretUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*openapi.APIClient)

	secretId, err := strconv.ParseInt(d.Id(), 10, 64)

	if err != nil {
		return diag.FromErr(err)
	}

	updateSecretRequest := openapi.NewUpdateSecretRequest()
	anyChange := false

	if d.HasChange("name") {
		secretName := d.Get("name").(string)
		updateSecretRequest.Name = &secretName
		anyChange = true
	}

	if d.HasChange("value") {
		secretValue := d.Get("value").(string)
		updateSecretRequest.Value = &secretValue
		anyChange = true
	}

	if anyChange {
		_, httpResp, err := client.SecretsApi.
			UpdateSecret(context.Background(), secretId).
			XRequestId(uuid.NewV4().String()).
			UpdateSecretRequest(*updateSecretRequest).
			Execute()

		if err != nil {
			return HandleResponseErrors(diags, httpResp)
		}

		d.Set("updated_at", time.Now().Format(time.RFC850))
		return resourceSecretRead(ctx, d, m)
	}

	return diags
}

func resourceSecretDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*openapi.APIClient)

	secretId, err := strconv.ParseInt(d.Id(), 10, 64)

	if err != nil {
		return diag.FromErr(err)
	}

	httpResp, err := client.SecretsApi.
		DeleteSecret(ctx, secretId).
		XRequestId(uuid.NewV4().String()).
		Execute()

	if err != nil {
		return HandleResponseErrors(diags, httpResp)
	}

	d.SetId("")

	return diags
}

func AddSecretToData(
	secret openapi.SecretResponse,
	d *schema.ResourceData,
	diags diag.Diagnostics,
) diag.Diagnostics {
	id := strconv.Itoa(int(secret.SecretId))
	if err := d.Set("id", id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", secret.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", secret.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("value", secret.Value); err != nil {
		return diag.FromErr(err)
	}
	createdAt := secret.CreatedAt.Format(time.RFC850)
	if err := d.Set("created_at", createdAt); err != nil {
		return diag.FromErr(err)
	}
	return diags
}
