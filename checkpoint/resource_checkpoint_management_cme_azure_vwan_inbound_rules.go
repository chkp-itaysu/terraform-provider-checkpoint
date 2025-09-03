package checkpoint

import (
	"encoding/json"
	"fmt"
	"log"

	checkpoint "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "A list of rules in a string format.",
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
		return fmt.Errorf("%s", err.Error())
	}

	inboundRules := AzureVwanInboundRulesRes.GetData()
	if checkIfRequestFailed(inboundRules) {
		errMessage := buildErrorMessage(inboundRules)
		return fmt.Errorf("%s", errMessage)
	}

	requestId := inboundRules["result"].(map[string]interface{})["request-id"].(string)

	requestErr := cmeWaitForReuqest(client, requestId)
	if requestErr != nil {
		return requestErr
	}

	d.SetId("")

	return nil
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
		rules := make([]map[string]interface{}, 0)
		if err := json.Unmarshal([]byte(v.(string)), &rules); err != nil {
			return fmt.Errorf("Failed to parse rules: %v. For more information on how a rule object should look like visit: <URL>", err)
		}

		payload["rules"] = rules

		log.Println("Parsed rules: ", payload["rules"])
	}

	log.Println("Create cme Azure VWAN inbound rules - NVA = ", nvaName)

	url := CmeApiPath + "/azure/virtualWANs/accounts/" + accountID + "/resourceGroups/" + nvaResourceGroup + "/inboundRules/" + nvaName

	AzureVwanInboundRulesRes, err := client.ApiCall(url, payload, client.GetSessionID(), true, client.IsProxyUsed())

	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	inboundRules := AzureVwanInboundRulesRes.GetData()
	if checkIfRequestFailed(inboundRules) {
		errMessage := buildErrorMessage(inboundRules)
		return fmt.Errorf("%s", errMessage)
	}

	requestId := inboundRules["result"].(map[string]interface{})["request-id"].(string)

	requestErr := cmeWaitForReuqest(client, requestId)
	if requestErr != nil {
		return requestErr
	}

	d.SetId("cme-azure-vwan-inbound-rules-" + nvaName + "-" + acctest.RandString(10))

	return dataSourceManagementCMEAzureVwanInboundRulesRead(d, m)
}
