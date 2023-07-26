package fybe

import (
	"context"
	"strconv"

	apiClient "fybe.com/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func dataSourceVPC() *schema.Resource {
	return &schema.Resource{
		Description: "Virtual private network contain your compute instances whereby they are able to communicate with each other in full isolation, using configured private IP addresses.",
		ReadContext: dataSourcePrivateNetworkRead,
		Schema: map[string]*schema.Schema{
			"created_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "VPC creation date.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "VPC modified date.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Id of this VPC.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The VPC name can consist of letters, numbers, colons, dashes, and underscores. However, it must not exceed 255 characters in length.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The VPC description should not exceed 255 characters in length.",
			},
			"instance_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Description: "Add compute instace Ids to the VPC. If you do not add any instance Ids an empty VPC will be created.",
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "compute instance Id.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Adjustable compute instance name, it is changeable in the cockpit.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Fybe name of the compute instance.",
						},
						"private_ip_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of all private IP addresses of the compute instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"v4": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "IP Address",
												},
												"netmask_cidr": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Netmask CIDR",
												},
												"gateway": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Gateway",
												},
											},
										},
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the VPC can be one of `ok|restart|reinstall|reinstallation failed|installing`",
						},
						"error_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If the instance is in an error state (see status property), the error message can be seen in this field.",
						},
					},
				},
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "us-east-1",
				Description: "Region slug of this VPC.",
			},
			"region_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Full name of the VPC region.",
			},
			"data_center": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specific data center where the Private Network is located.",
			},
			"available_ips": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total available IPs in the VPC.",
			},
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cidr range of the VPC.",
			},
		},
	}
}

func dataSourcePrivateNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*apiClient.APIClient)

	var privateNetworktId int64
	var err error
	id := d.Get("id").(string)
	if id != "" {
		privateNetworktId, err = strconv.ParseInt(id, 10, 64)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	res, httpResp, err := client.VirtualPrivateCloudVPCApi.
		RetrievePrivateNetwork(ctx, privateNetworktId).
		XRequestId(uuid.NewV4().String()).
		Execute()

	if err != nil {
		return HandleResponseErrors(diags, httpResp)
	} else if len(res.Data) != 1 {
		return MultipleDataObjectsError(diags)
	}

	d.SetId(strconv.Itoa(int(res.Data[0].PrivateNetworkId)))

	return AddPrivateNetworkToData(res.Data[0], d, diags)
}
