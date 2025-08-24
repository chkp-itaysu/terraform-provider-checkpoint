package checkpoint

import (
	"fmt"
	checkpoint "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceManagementCMEAzureVwanInboundRule() *schema.Resource {
	return &schema.Resource{
		Create: createManagementCMEAzureVwanInboundRule,
		Update: createManagementCMEAzureVwanInboundRule,
		Read:   dataSourceManagementCMEAzureVwanInboundRulesRead,
		Delete: deleteManagementCMEAzureVwanInboundRule,
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
	}
}

func deleteManagementCMEAzureVwanInboundRule(d *schema.ResourceData, m interface{}) error {
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

	d.SetId("")

	return nil
}

func createManagementCMEAzureVwanInboundRule(d *schema.ResourceData, m interface{}) error {
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
	ruleMap := make(map[string]interface{}, 0)

	if v, ok := d.GetOk("name"); ok {
		ruleMap["name"] = v.(string)
	}

	if v, ok := d.GetOk("lb_public_ips"); ok {
		ruleMap["lb_public_ips"] = v.([]interface{})
	}

	if v, ok := d.GetOk("original_ports"); ok {
		ruleMap["original_ports"] = v.([]interface{})
	}

	if v, ok := d.GetOk("original_source"); ok {
		ruleMap["original_source"] = v.(string)
	}

	if v, ok := d.GetOk("protocol"); ok {
		ruleMap["protocol"] = v.(string)
	}

	payload["rules"] = []map[string]interface{}{ruleMap}

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

	d.SetId("cme-azure-vwan-inbound-rules-" + ruleMap["name"].(string) + acctest.RandString(10))

	return dataSourceManagementCMEAzureVwanInboundRulesRead(d, m)
}
