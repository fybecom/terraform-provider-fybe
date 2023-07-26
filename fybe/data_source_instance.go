package fybe

import (
	"context"
	"strconv"

	apiClient "fybe.com/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Create and manage compute resources. Cloud-Init is also supported.",
		ReadContext: dataSourceInstanceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Compute instance Id",
			},
			"last_updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified date of the instance.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Fybe name of the instance.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Your choosed name for the instance, you can change it in the cockpit.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CAUTION: On updating this value your server will be reinstalled! You can find the available image Ids via [API](https://api.fybe.com/#tag/Images/operation/retrieveImage) or via our [command line](https://github.com/fybecom/fybe) tool with this command: `fybe get images`.",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in it the instance should be installed. Default region is the us-east-1.",
			},
			"product_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Choose the instance that fits for you.",
			},
			"ip_config": {
				Type:     schema.TypeList,
				Computed: true,
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
						"v6": {
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
			"mac_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance mac address.",
			},
			"ram_mb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Image ram size in megabyte.",
			},
			"cpu_cores": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "CPU core count of the instance.",
			},
			"disk_mb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Image disk size of the instance in megabyte.",
			},
			"os_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of operating system (OS) installed on the instance.",
			},
			"ssh_keys": {
				Computed: true,
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Important: Updating this value will result in your server being reinstalled! This pertains to an array of secretIds that represent public SSH keys used for logging in as the defaultUser with administrator/root privileges. This functionality is applicable to Linux/BSD systems. For further details, please consult the Secrets Management API.",
			},
			"created_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance creation date.",
			},
			"cancel_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The date on which the instance will not longer be accessable due to a cancellation.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance status. The status can be set to `provisioning|uninstalled|running|stopped|error|installing|unknown|installed`.",
			},
			"v_host_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "vHost Identifier.",
			},
			"add_ons": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "Ids for addons to pimp your instance.",
						},
						"quantity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The number of Addons you wish to buy.",
						},
					},
				},
			},
			"error_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If the instance is in an error state (see status property), the error message can be seen in this field.",
			},
			"product_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instances category depending on ProductId. Following product types are available: `ssd`,`nvme`.",
			},
			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "CAUTION: On updating this value your server will be reinstalled! Cloud-Init Config in order to customize during start of the instance.",
			},
			"license": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Additional license in order to enhance your instance. Possible license to buy are cpanel and plesk.",
			},
			"period": {
				Type:        schema.TypeInt,
				Computed:    true,
				Default:     1,
				Description: "",
			},
			"additional_ips_v4": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All other additional IP addresses of the instance.",
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
	}
}

func dataSourceInstanceRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*apiClient.APIClient)

	var instanceId int64
	var err error
	id := d.Get("id").(string)
	if id != "" {
		instanceId, err = strconv.ParseInt(id, 10, 64)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	res, httpResp, err := client.InstancesApi.
		RetrieveInstance(ctx, int64(instanceId)).
		XRequestId(uuid.NewV4().String()).
		Execute()

	if err != nil {
		return HandleResponseErrors(diags, httpResp)
	} else if len(res.Data) != 1 {
		return MultipleDataObjectsError(diags)
	}

	d.SetId(strconv.Itoa(int(res.Data[0].InstanceId)))

	return AddInstanceToData(
		res.Data[0],
		d,
		diags,
	)
}
