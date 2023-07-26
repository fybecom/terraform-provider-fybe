package fybe

import (
	"context"

	apiClient "fybe.com/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func dataSourceImage() *schema.Resource {
	return &schema.Resource{
		Description: "To add a custom image, provide a direct download URL for the image in either .iso or .qcow2 format. Any other formats will not be accepted. Please be aware that the download duration depends on your network speed and the image's size. You can monitor the download status by making a GET request to obtain information about the image. Download requests will be declined if they surpass your allocated limits.",
		ReadContext: dataSourceImageRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The identifier of the image",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time of the image.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Image name.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Description of the image.",
			},
			"uploaded_size_mb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the uploaded image in megabyte.",
			},
			"os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Please indicate the type of operating system (OS) you want to use. For MS Windows, specify 'Windows,' and for other OS options, specify 'Linux.' Selecting the incorrect OS type may result in a non-functional cloud instance.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Please provide a version number to differentiate the contents of the image, such as the version of the operating system.",
			},
			"format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Format of your image `iso` or `qcow`.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Downloading status of the image `downloading|downloaded|error`.",
			},
			"error_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If the image is in an error state (see status property), the error message can be seen in this field.",
			},
			"standard_image": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag indicating that the image is either a standard (true) or a custom image (false).",
			},
			"creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation date of the image.",
			},
		},
	}
}

func dataSourceImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*apiClient.APIClient)

	imageId := d.Get("id").(string)

	if imageId == "" {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "imageId should not be empty",
		})
	}

	res, httpResp, err := client.ImagesApi.
		RetrieveImage(ctx, imageId).
		XRequestId(uuid.NewV4().String()).
		Execute()

	if err != nil {
		return HandleResponseErrors(diags, httpResp)
	} else if len(res.Data) != 1 {
		return MultipleDataObjectsError(diags)
	}

	d.SetId(res.Data[0].ImageId)

	return AddImageToData(res.Data[0], d, diags)
}
