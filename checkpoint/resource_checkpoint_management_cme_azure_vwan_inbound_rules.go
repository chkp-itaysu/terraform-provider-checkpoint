package checkpoint

import (
	"fmt"
	checkpoint "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceManagementCMEAzureVwanInboundRules() *schema.Resource {
	return &schema.Resource{
		Create: createManagementCMEAzureVwanInboundRules,
		Update: createManagementCMEAzureVwanInboundRules,
		Read:   dataSourceManagementCMEAzureVwanInboundRulesRead,
		Delete: deleteManagementCMEAzureVwanInboundRules,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
				Required:    true,
				Description: "A list of inbound rules of the NVA.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The rule name.",
						},
						"lb_public_ips": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "A list of outbound public IPs for that rule.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"original_ports": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "The list of ports allowed in that rule.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"original_source": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The rule inbound IP address.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The traffic protocol in that rule.",
						},
					},
				},
			},
		},
	}
}

func deleteManagementCMEAzureVwanInboundRules(d *schema.ResourceData, m interface{}) error {
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

	log.Println("Delete cme Azure VWAN inbound rules - NVA = ", nvaName)

	url := CmeApiPath + "/azure/virtualWANs/accounts/" + accountID + "/resourceGroups/" + nvaResourceGroup + "/inboundRules/" + nvaName

	AzureVwanInboundRulesRes, err := client.ApiCall(url, nil, client.GetSessionID(), true, client.IsProxyUsed(), "DELETE")

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	inboundRules := AzureVwanInboundRulesRes.GetData()
	if checkIfRequestFailed(inboundRules) {
		errMessage := buildErrorMessage(inboundRules)
		return fmt.Errorf(errMessage)
	}

	requestId := inboundRules["result"].(map[string]interface{})["request-id"].(string)

	requestErr := cmeWaitForReuqest(client, requestId)
	if requestErr != nil {
		return requestErr
	}

	d.SetId("cme-azure-vwan-inbound-rules-" + accountID + "-" + nvaResourceGroup + "-" + nvaName + "-" + acctest.RandString(10))

	return dataSourceManagementCMEAzureVwanInboundRulesRead(d, m)
}

func createManagementCMEAzureVwanInboundRules(d *schema.ResourceData, m interface{}) error {
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

	payload := make(map[string]interface{}, 0)

	if v, ok := d.GetOk("rules"); ok {
		rules := v.([]interface{})
		rulesList := make([]map[string]interface{}, 0, len(rules))
		for _, rule := range rules {
			ruleMap := rule.(map[string]interface{})
			rulePayload := make(map[string]interface{})
			rulePayload["name"] = ruleMap["name"].(string)
			rulePayload["lb_public_ips"] = ruleMap["lb_public_ips"].([]interface{})
			rulePayload["original_ports"] = ruleMap["original_ports"].([]interface{})
			rulePayload["original_source"] = ruleMap["original_source"].(string)
			rulePayload["protocol"] = ruleMap["protocol"].(string)
			rulesList = append(rulesList, rulePayload)
		}
		payload["rules"] = rulesList
	} else {
		return fmt.Errorf("rules must be provided")
	}

	log.Println("Create cme Azure VWAN inbound rules - NVA = ", nvaName)

	url := CmeApiPath + "/azure/virtualWANs/accounts/" + accountID + "/resourceGroups/" + nvaResourceGroup + "/inboundRules/" + nvaName

	AzureVwanInboundRulesRes, err := client.ApiCall(url, payload, client.GetSessionID(), true, client.IsProxyUsed())

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	inboundRules := AzureVwanInboundRulesRes.GetData()
	if checkIfRequestFailed(inboundRules) {
		errMessage := buildErrorMessage(inboundRules)
		return fmt.Errorf(errMessage)
	}

	requestId := inboundRules["result"].(map[string]interface{})["request-id"].(string)

	requestErr := cmeWaitForReuqest(client, requestId)
	if requestErr != nil {
		return requestErr
	}

	d.SetId("cme-azure-vwan-inbound-rules-" + accountID + "-" + nvaResourceGroup + "-" + nvaName + "-" + acctest.RandString(10))

	return dataSourceManagementCMEAzureVwanInboundRulesRead(d, m)
}
