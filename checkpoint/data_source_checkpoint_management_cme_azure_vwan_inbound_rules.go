package checkpoint

import (
	"fmt"
	checkpoint "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func dataSourceManagementCMEAzureVwanInboundRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceManagementCMEAzureVwanInboundRulesRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the Azure account.",
			},
			"nva_resource_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the resource group that contains the NVA.",
			},
			"nva_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the NVA.",
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of inbound rules of the NVA.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The rule name.",
						},
						"lb_public_ips": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of outbound public IPs for that rule.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"original_ports": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of ports allowed in that rule.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"original_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The rule inbound IP address.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The traffic protocol in that rule.",
						},
					},
				},
			},
		},
	}
}

func dataSourceManagementCMEAzureVwanInboundRulesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*checkpoint.ApiClient)

	var accountID string
	var nvaResourceGroup string
	var nvaName string

	if v, ok := d.GetOk("account_id"); ok {
		accountID = v.(string)
	}
	if v, ok := d.GetOk("nva_resource_group"); ok {
		nvaResourceGroup = v.(string)
	}
	if v, ok := d.GetOk("nva_name"); ok {
		nvaName = v.(string)
	}

	log.Println("Read cme Azure VWAN inbound rules - NVA = ", nvaName)

	url := CmeApiPath + "/azure/virtualWANs/accounts/" + accountID + "/resourceGroups/" + nvaResourceGroup + "/inboundRules/" + nvaName

	AzureVwanInboundRulesRes, err := client.ApiCall(url, nil, client.GetSessionID(), true, client.IsProxyUsed(), "GET")

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	inboundRules := AzureVwanInboundRulesRes.GetData()
	if checkIfRequestFailed(inboundRules) {
		errMessage := buildErrorMessage(inboundRules)
		return fmt.Errorf(errMessage)
	}

	if v, ok := d.GetOk("id"); !ok || v.(string) == "" {
		d.SetId("cme-azure-vwan-inbound-rules-" + accountID + "-" + nvaResourceGroup + "-" + nvaName + "-" + acctest.RandString(10))
	}

	rulesList := inboundRules["result"].(map[string]interface{})["rules"].([]interface{})
	var rulesListToReturn []map[string]interface{}

	for i := range rulesList {
		rulesListToReturn = append(rulesListToReturn, rulesList[i].(map[string]interface{}))
	}
	_ = d.Set("rules", rulesListToReturn)

	log.Println("Inbound Rules: ", d.Get("rules"))

	return nil
}
