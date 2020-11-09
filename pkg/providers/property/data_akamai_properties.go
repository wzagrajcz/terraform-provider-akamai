package property

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/akamai"
	"github.com/akamai/terraform-provider-akamai/v2/pkg/tools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAkamaiProperties() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAkamaiPropertiesRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contract_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"properties": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "List of properties",
			},
		},
	}
}

func dataAkamaiPropertiesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	meta := akamai.Meta(m)
	log := meta.Log("PAPI", "dataAkamaiPropertiesRead")
	// create a context with logging for api calls
	ctx = session.ContextWithOptions(
		ctx,
		session.WithContextLog(log),
	)
	log.Debug("Listing Properties")

	groupId, err := tools.GetStringValue("group_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	contractId, err := tools.GetStringValue("contract_id", d)
	if err != nil {
		return diag.FromErr(err)
	}
	groupId = tools.AddPrefix(groupId, "grp_")
	contractId = tools.AddPrefix(contractId, "ctr_")
	properties, err := getProperties(ctx, groupId, contractId, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing properties: %w", err))
	}
	// setting concatenated id to uniquely identify data
	d.SetId(groupId + contractId)
	props, err := json.Marshal(properties)
	if err != nil {
		return diag.FromErr(err)
	}
	/* setting raw json here in current scope. We have to set all the json fields in
	property struct for more granular access to properties object */
	if err := d.Set("properties", string(props)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting properties: %s", err))
	}
	return nil
}

// Reusable function to fetch all the properties for a given group and contract
func getProperties(ctx context.Context, groupId string, contractId string, meta akamai.OperationMeta) (*papi.GetPropertiesResponse, error) {
	client := inst.Client(meta)
	req := papi.GetPropertiesRequest{
		ContractID: contractId,
		GroupID:    groupId,
	}
	props, err := client.GetProperties(ctx, req)
	if err != nil {
		return nil, err
	}
	return props, nil
}
