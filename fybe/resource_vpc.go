package fybe

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	openapi "fybe.com/apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	uuid "github.com/satori/go.uuid"
)

var httpConflict string = "409 Conflict"

func resourceVPC() *schema.Resource {
	return &schema.Resource{
		Description:   "Virtual private network contain your compute instances whereby they are able to communicate with each other in full isolation, using configured private IP addresses.",
		CreateContext: resourcePrivateNetworkCreate,
		ReadContext:   resourcePrivateNetworkRead,
		UpdateContext: resourcePrivateNetworkUpdate,
		DeleteContext: resourcePrivateNetworkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
							Type:        schema.TypeInt,
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

func resourcePrivateNetworkCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*openapi.APIClient)

	privateNetworkName := d.Get("name").(string)
	privateNetworkDescription := d.Get("description").(string)
	privateNetworkRegion := d.Get("region").(string)

	createPrivateNetworkRequest := openapi.NewCreatePrivateNetworkRequestWithDefaults()
	createPrivateNetworkRequest.Name = privateNetworkName
	createPrivateNetworkRequest.Description = &privateNetworkDescription
	createPrivateNetworkRequest.Region = &privateNetworkRegion

	res, httpResp, err := client.VirtualPrivateCloudVPCApi.
		CreatePrivateNetwork(context.Background()).
		XRequestId(uuid.NewV4().String()).
		CreatePrivateNetworkRequest(*createPrivateNetworkRequest).
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
	instancesToAdd := d.Get("instance_ids").(*schema.Set).List()
	privateNetworkId := res.Data[0].PrivateNetworkId

	for _, instanceId := range instancesToAdd {
		instanceIdInt := instanceId.(int)
		instanceId := int64(instanceIdInt)

		if err != nil && !strings.Contains(err.Error(), httpConflict) {
			return HandleResponseErrors(diags, httpResp)
		}
		httpResp, err = assignInstanceToPrivateNetwork(diags, client, privateNetworkId, instanceId)
		if err != nil {
			return HandleResponseErrors(diags, httpResp)
		}
	}
	d.SetId(strconv.Itoa(int(privateNetworkId)))
	return resourcePrivateNetworkRead(ctx, d, m)
}

func assignInstanceToPrivateNetwork(
	diags diag.Diagnostics,
	client *openapi.APIClient,
	privateNetworkId,
	instanceId int64) (*http.Response, error) {

	pollInstance(diags, client, instanceId)

	_, httpResp, err := client.VirtualPrivateCloudVPCApi.AssignInstancePrivateNetwork(
		context.Background(),
		privateNetworkId,
		instanceId).XRequestId(uuid.NewV4().String()).Execute()

	return httpResp, err
}

func unassignInstanceToPrivateNetwork(
	diags diag.Diagnostics,
	client *openapi.APIClient,
	privateNetworkId int64,
	instanceId int64) (*http.Response, error) {

	_, httpResp, err := client.VirtualPrivateCloudVPCApi.UnassignInstancePrivateNetwork(
		context.Background(),
		privateNetworkId,
		instanceId).XRequestId(uuid.NewV4().String()).Execute()

	return httpResp, err
}

func resourcePrivateNetworkRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*openapi.APIClient)

	privateNetworkId, err := strconv.ParseInt(d.Id(), 10, 64)

	if err != nil {
		return diag.FromErr(err)
	}

	res, httpResp, err := client.VirtualPrivateCloudVPCApi.
		RetrievePrivateNetwork(ctx, privateNetworkId).
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

	return AddPrivateNetworkToData(res.Data[0], d, diags)
}

func resourcePrivateNetworkUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*openapi.APIClient)

	privateNetworkId, err := strconv.ParseInt(d.Id(), 10, 64)

	if err != nil {
		return diag.FromErr(err)
	}

	updatePrivateNetworkRequest := openapi.NewPatchPrivateNetworkRequest()
	anyChange := false

	if d.HasChange("name") {
		privateNetworkName := d.Get("name").(string)
		updatePrivateNetworkRequest.Name = &privateNetworkName
		anyChange = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updatePrivateNetworkRequest.Description = &description
		anyChange = true
	}

	if d.HasChange("instance_ids") {
		rsltDiag := handlePrivateNetworkInstanceChanges(diags, d, client, privateNetworkId)
		if rsltDiag != nil {
			return rsltDiag
		}
		anyChange = true
	}

	if anyChange {
		_, httpResp, err := client.VirtualPrivateCloudVPCApi.
			PatchPrivateNetwork(context.Background(), privateNetworkId).
			XRequestId(uuid.NewV4().String()).
			PatchPrivateNetworkRequest(*updatePrivateNetworkRequest).
			Execute()

		if err != nil {
			return HandleResponseErrors(diags, httpResp)
		}

		d.Set("updated_at", time.Now().Format(time.RFC850))
		return resourcePrivateNetworkRead(ctx, d, m)
	}
	return diags
}

func handlePrivateNetworkInstanceChanges(diags diag.Diagnostics,
	d *schema.ResourceData,
	client *openapi.APIClient,
	privateNetworkId int64) diag.Diagnostics {

	//Remove instances which are not more in this private network
	old, new := d.GetChange("instance_ids")
	oldInstanceIds := old.(*schema.Set).List()
	for _, instanceId := range oldInstanceIds {
		instanceIdInt := instanceId.(int)
		instanceId := int64(instanceIdInt)

		httpResp, err := unassignInstanceToPrivateNetwork(diags, client, privateNetworkId, instanceId)
		if err != nil {
			return HandleResponseErrors(diags, httpResp)
		}
	}

	//Add new instances which are now in this private network
	newInstanceIds := new.(*schema.Set).List()
	for _, instanceId := range newInstanceIds {
		instanceIdInt := instanceId.(int)
		instanceId := int64(instanceIdInt)

		httpResp, err := assignInstanceToPrivateNetwork(diags, client, privateNetworkId, instanceId)
		if err != nil {
			return HandleResponseErrors(diags, httpResp)
		}
	}
	return nil
}

func resourcePrivateNetworkDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*openapi.APIClient)

	privateNetworkId, err := strconv.ParseInt(d.Id(), 10, 64)

	if err != nil {
		return diag.FromErr(err)
	}

	readRes, httpResp, err := client.VirtualPrivateCloudVPCApi.
		RetrievePrivateNetwork(ctx, privateNetworkId).
		XRequestId(uuid.NewV4().String()).
		Execute()

	if err != nil {
		return HandleResponseErrors(diags, httpResp)
	}

	for _, i := range readRes.Data[0].Instances {
		client.VirtualPrivateCloudVPCApi.UnassignInstancePrivateNetwork(ctx, privateNetworkId, i.InstanceId).XRequestId(uuid.NewV4().String()).Execute()
	}

	httpResp, err = client.VirtualPrivateCloudVPCApi.
		DeletePrivateNetwork(ctx, privateNetworkId).
		XRequestId(uuid.NewV4().String()).
		Execute()

	if err != nil {
		return HandleResponseErrors(diags, httpResp)
	}

	d.SetId("")

	return diags
}

func AddPrivateNetworkToData(
	privateNetwork openapi.PrivateNetworkResponse,
	d *schema.ResourceData,
	diags diag.Diagnostics,
) diag.Diagnostics {
	id := strconv.Itoa(int(privateNetwork.PrivateNetworkId))
	if err := d.Set("id", id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", privateNetwork.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", privateNetwork.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("region", privateNetwork.Region); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data_center", privateNetwork.DataCenter); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("available_ips", privateNetwork.AvailableIps); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cidr", privateNetwork.Cidr); err != nil {
		return diag.FromErr(err)
	}
	createdDate := privateNetwork.CreatedDate.Format(time.RFC850)
	if err := d.Set("created_date", createdDate); err != nil {
		return diag.FromErr(err)
	}

	instanceIds := []int64{}
	instances := []map[string]interface{}{}

	for _, instance := range privateNetwork.Instances {
		instanceIds = append(instanceIds, instance.InstanceId)
		instances = append(instances, buildInstanceIpConfig(instance))
	}
	if err := d.Set("instance_ids", instanceIds); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("instances", instances); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func buildInstanceIpConfig(instance openapi.Instances) map[string]interface{} {
	instanceConfig := make(map[string]interface{})

	instanceConfig["instance_id"] = instance.InstanceId
	instanceConfig["display_name"] = instance.DisplayName
	instanceConfig["name"] = instance.Name
	instanceConfig["status"] = instance.Status
	instanceConfig["error_message"] = instance.ErrorMessage

	privateIpConfig := make(map[string]interface{})
	privateIpConfigList := []map[string]interface{}{}

	for _, privateIpConfigV4 := range instance.PrivateIpConfig.V4 {
		ipConfig := make(map[string]interface{})
		ipConfig["ip"] = privateIpConfigV4.Ip
		ipConfig["netmask_cidr"] = privateIpConfigV4.NetmaskCidr
		ipConfig["gateway"] = privateIpConfigV4.Gateway
		privateIpConfigList = append(privateIpConfigList, ipConfig)
	}

	privateIpConfig["v4"] = privateIpConfigList
	instanceConfig["private_ip_config"] = []interface{}{privateIpConfig}

	return instanceConfig
}
