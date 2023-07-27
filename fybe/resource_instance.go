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

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Description:   "Create and manage compute resources. Cloud-Init is also supported.",
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
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
			"root_password": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "CAUTION: On updating this value your server will be reinstalled! Root password of the instance.",
			},
			"created_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance creation date.",
			},
			"cancel_date": {
				Type:        schema.TypeString,
				Optional:    true,
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
				Optional:    true,
				Computed:    true,
				Description: "",
			},
			"additional_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "All other additional IP addresses of the instance.",
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
		},
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*openapi.APIClient)

	createInstanceRequest := openapi.NewCreateInstanceRequestWithDefaults()

	displayName := d.Get("display_name").(string)
	imageId := d.Get("image_id").(string)
	region := d.Get("region").(string)
	productId := d.Get("product_id").(string)
	sshKeys := d.Get("ssh_keys")
	rootPassword := d.Get("root_password")
	userData := d.Get("user_data").(string)
	license := d.Get("license").(string)
	period := d.Get("period").(int)

	if displayName != "" {
		createInstanceRequest.DisplayName = &displayName
	}
	if imageId != "" {
		createInstanceRequest.ImageId = &imageId
	}
	if region != "" {
		createInstanceRequest.Region = &region
	}
	if productId != "" {
		createInstanceRequest.ProductId = productId
	}
	if sshKeys != nil {
		var sshKeys64 []int64
		for _, key := range sshKeys.([]interface{}) {
			sshKey := key.(int)
			sshKeys64 = append(sshKeys64, int64(sshKey))
		}
		createInstanceRequest.SshKeys = &sshKeys64
	}
	if rootPassword != nil {
		rootPassword64 := int64(rootPassword.(int))
		createInstanceRequest.RootPassword = &rootPassword64
	}
	if userData != "" {
		createInstanceRequest.UserData = &userData
	}
	if license != "" {
		createInstanceRequest.License = &license
	}
	if period != 0 {
		createInstanceRequest.Period = int64(period)
	}

	res, httpResp, err := client.InstancesApi.
		CreateInstance(ctx).
		XRequestId(uuid.NewV4().String()).
		CreateInstanceRequest(*createInstanceRequest).
		Execute()

	if err != nil {
		return HandleResponseErrors(diags, httpResp)
	} else if len(res.Data) != 1 {
		return MultipleDataObjectsError(diags)
	}

	d.SetId(strconv.Itoa(int(res.Data[0].InstanceId)))

	return resourceInstanceRead(ctx, d, m)
}

func resourceInstanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*openapi.APIClient)

	instanceId, err := strconv.ParseInt(d.Id(), 10, 64)

	if err != nil {
		return diag.FromErr(err)
	}

	instance, diag := pollInstanceInstalled(diags, client, ctx, instanceId)

	if err != nil || instance == nil {
		return append(diags, diag...)
	}

	return AddInstanceToData(*instance, d, diags)
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*openapi.APIClient)
	instanceId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	if shouldUpdateInstanceValues(d) {
		updateInstanceValues(d, client, ctx, instanceId, diags, m)
	}
	if shouldReinstall(d) {
		reinstall(d, client, ctx, instanceId, diags, m)
	}
	return resourceInstanceRead(ctx, d, m)
}

func shouldUpdateInstanceValues(d *schema.ResourceData) bool {
	updateChange := false
	if d.HasChange("display_name") {
		displayname := d.Get("display_name").(string)
		if displayname != "" {
			updateChange = true
		}
	}
	return updateChange
}

func shouldReinstall(d *schema.ResourceData) bool {
	reinstallChange := false
	if d.HasChange("ssh_keys") {
		sshKeys := d.Get("ssh_keys")
		if sshKeys != nil {
			reinstallChange = true
		}
	}

	if d.HasChange("root_password") {
		rootPassword := d.Get("root_password")
		if rootPassword != nil {
			reinstallChange = true
		}
	}

	if d.HasChange("user_data") {
		userData := d.Get("user_data").(string)
		if userData != "" {
			reinstallChange = true
		}
	}
	if d.HasChange("image_id") {
		imageId := d.Get("image_id").(string)
		if imageId != "" {
			reinstallChange = true
		}
	}
	return reinstallChange
}

func updateInstanceValues(d *schema.ResourceData, client *openapi.APIClient, ctx context.Context, instanceId int64, diags diag.Diagnostics, m interface{}) diag.Diagnostics {

	patchInstanceRequest := *openapi.NewPatchInstanceRequestWithDefaults()
	displayName := d.Get("display_name").(string)
	patchInstanceRequest.DisplayName = &displayName

	res, httpResp, err := client.InstancesApi.
		PatchInstance(context.Background(), instanceId).
		PatchInstanceRequest(patchInstanceRequest).
		XRequestId(uuid.NewV4().String()).Execute()

	if err != nil {
		return HandleResponseErrors(diags, httpResp)
	} else if len(res.Data) != 1 {
		return MultipleDataObjectsError(diags)
	}
	d.SetId(strconv.Itoa(int(res.Data[0].InstanceId)))
	return diags
}

func reinstall(d *schema.ResourceData, client *openapi.APIClient, ctx context.Context, instanceId int64, diags diag.Diagnostics, m interface{}) diag.Diagnostics {
	patchInstanceRequest := openapi.NewReinstallInstanceRequestWithDefaults()

	if d.HasChange("ssh_keys") {
		sshKeys := d.Get("ssh_keys")
		if sshKeys != nil {
			var sshKeys64 []int64
			for _, key := range sshKeys.([]interface{}) {
				sshKey := key.(int)
				sshKeys64 = append(sshKeys64, int64(sshKey))
			}
			patchInstanceRequest.SshKeys = &sshKeys64
		}
	}

	if d.HasChange("root_password") {
		rootPassword := d.Get("root_password")
		if rootPassword != nil {
			rootPassword64 := int64(rootPassword.(int))
			patchInstanceRequest.RootPassword = &rootPassword64
		}
	}

	if d.HasChange("user_data") {
		userData := d.Get("user_data").(string)
		if userData != "" {
			patchInstanceRequest.UserData = &userData
		}
	}

	imageId := d.Get("image_id").(string)
	if imageId != "" {
		patchInstanceRequest.ImageId = imageId
	}

	res, httpResp, err := client.InstancesApi.
		ReinstallInstance(ctx, instanceId).
		XRequestId(uuid.NewV4().String()).
		ReinstallInstanceRequest(*patchInstanceRequest).
		Execute()

	if err != nil {
		return HandleResponseErrors(diags, httpResp)
	} else if len(res.Data) != 1 {
		return MultipleDataObjectsError(diags)
	}
	d.SetId(strconv.Itoa(int(res.Data[0].InstanceId)))
	return diags
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func AddInstanceToData(
	instance openapi.InstanceResponse,
	d *schema.ResourceData,
	diags diag.Diagnostics,
) diag.Diagnostics {
	id := strconv.Itoa(int(instance.InstanceId))
	if err := d.Set("id", id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", instance.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("display_name", instance.DisplayName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("image_id", instance.ImageId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("product_id", instance.ProductId); err != nil {
		return diag.FromErr(err)
	}
	ipConfig := buildIpConfig(&instance.IpConfig)
	if err := d.Set("ip_config", ipConfig); err != nil && len(ipConfig) > 0 {
		return diag.FromErr(err)
	}
	if err := d.Set("mac_address", instance.MacAddress); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ram_mb", instance.RamMb); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cpu_cores", instance.CpuCores); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("disk_mb", instance.DiskMb); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("os_type", instance.OsType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ssh_keys", instance.SshKeys); err != nil {
		return diag.FromErr(err)
	}
	createdDate := instance.CreatedDate.Format(time.RFC850)
	if err := d.Set("created_date", createdDate); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cancel_date", instance.CancelDate); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("status", instance.Status); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("v_host_id", instance.VHostId); err != nil {
		return diag.FromErr(err)
	}
	addOns := buildAddons(instance.AddOns)
	if err := d.Set("add_ons", addOns); err != nil && len(addOns) > 0 {
		return diag.FromErr(err)
	}
	if err := d.Set("error_message", instance.ErrorMessage); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("product_type", instance.ProductType); err != nil {
		return diag.FromErr(err)
	}
	additionalIps := buildAdditionalIps(instance.AdditionalIps)
	if err := d.Set("additional_ips", additionalIps); err != nil &&
		len(additionalIps) > 0 {
		return diag.FromErr(err)
	}

	return diags
}

func buildIpConfig(ipConfigResponse *openapi.IpConfig1) []interface{} {
	if ipConfigResponse != nil {
		ipConfig := make(map[string]interface{})

		v4 := make(map[string]interface{})
		v4["ip"] = ipConfigResponse.V4.Ip
		v4["netmask_cidr"] = int(ipConfigResponse.V4.NetmaskCidr)
		v4["gateway"] = ipConfigResponse.V4.Gateway

		v6 := make(map[string]interface{})
		v6["ip"] = ipConfigResponse.V6.Ip
		v6["netmask_cidr"] = int(ipConfigResponse.V6.NetmaskCidr)
		v6["gateway"] = ipConfigResponse.V6.Gateway

		ipConfig["v4"] = []interface{}{v4}
		ipConfig["v6"] = []interface{}{v6}

		return []interface{}{ipConfig}
	}

	return nil
}

func buildAddons(addOnResponse []openapi.AddOnResponse) []map[string]interface{} {
	if addOnResponse != nil {
		var addOns []map[string]interface{}

		for _, addOn := range addOnResponse {
			builtAddOn := make(map[string]interface{})
			builtAddOn["id"] = strconv.FormatInt(addOn.Id, 10)
			builtAddOn["quantity"] = addOn.Quantity

			addOns = append(addOns, builtAddOn)
		}

		return addOns
	}
	return nil
}

func buildAdditionalIps(
	additionalIpsResponse []openapi.AdditionalIp,
) []map[string]interface{} {
	if additionalIpsResponse != nil {
		additionalIps := []map[string]interface{}{}

		for _, ipV4 := range additionalIpsResponse {
			v4 := make(map[string]interface{})
			ipConfig := make(map[string]interface{})
			ipConfig["ip"] = ipV4.V4.Ip
			ipConfig["netmask_cidr"] = ipV4.V4.NetmaskCidr
			ipConfig["gateway"] = ipV4.V4.Gateway

			v4["v4"] = []interface{}{ipConfig}
			additionalIps = append(additionalIps, v4)
		}

		return additionalIps
	}

	return nil
}

func pollInstanceInstalled(
	diags diag.Diagnostics,
	client *openapi.APIClient,
	ctx context.Context,
	instanceId int64,
) (*openapi.InstanceResponse, diag.Diagnostics) {
	res, httpResp, err := client.InstancesApi.
		RetrieveInstance(ctx, instanceId).
		XRequestId(uuid.NewV4().String()).
		Execute()

	if err != nil {
		return nil, HandleResponseErrors(diags, httpResp)
	} else if len(res.Data) != 1 {
		return nil, MultipleDataObjectsError(diags)
	}

	status := res.Data[0].Status

	if status == openapi.PROVISIONING || status == openapi.INSTALLING {
		time.Sleep(time.Second)
		return pollInstanceInstalled(diags, client, ctx, instanceId)
	}

	return &res.Data[0], nil
}
