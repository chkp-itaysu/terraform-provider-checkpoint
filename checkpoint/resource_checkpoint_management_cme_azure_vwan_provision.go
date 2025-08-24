package checkpoint

import (
	"fmt"
	checkpoint "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceManagementCMEAzureVwanProvision() *schema.Resource {
	return &schema.Resource{
		Create: createManagementCMEAzureVwanProvision,
		Update: updateManagementCMEAzureVwanProvision,
		Read:   readManagementCMEAzureVwanProvision,
		Delete: deleteManagementCMEAzureVwanProvision,
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
			"base64_sic_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The sic key in base64 format.",
			},
			"policy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The provision policy.",
			},
			"autonomous_threat_prevention": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable autonomous threat prevention.",
			},
			"identity_awareness": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable identity awareness.",
			},
		},
	}
}

func createManagementCMEAzureVwanProvision(d *schema.ResourceData, m interface{}) error {
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

	if v, ok := d.GetOk("base64_sic_key"); ok {
		payload["base64_sic_key"] = v.(string)
	}

	if v, ok := d.GetOk("policy"); ok {
		payload["policy"] = v.(string)
	}

	if v, ok := d.GetOk("autonomous_threat_prevention"); ok {
		payload["autonomous_threat_prevention"] = v.(bool)
	}

	if v, ok := d.GetOk("identity_awareness"); ok {
		payload["identity_awareness"] = v.(bool)
	}

	log.Println("Create cme Azure VWAN provisioning - NVA = ", nvaName)

	url := CmeApiPath + "/azure/virtualWANs/accounts/" + accountID + "/resourceGroups/" + nvaResourceGroup + "/provision/" + nvaName

	AzureVwanProvisionRes, err := client.ApiCall(url, payload, client.GetSessionID(), true, client.IsProxyUsed())

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	provision := AzureVwanProvisionRes.GetData()
	if checkIfRequestFailed(provision) {
		errMessage := buildErrorMessage(provision)
		return fmt.Errorf(errMessage)
	}

	requestId := provision["result"].(map[string]interface{})["request-id"].(string)

	requestErr := cmeWaitForReuqest(client, requestId)
	if requestErr != nil {
		return requestErr
	}

	d.SetId("cme-azure-vwan-provision-" + accountID + "-" + nvaResourceGroup + "-" + nvaName + "-" + acctest.RandString(10))

	return nil
}

func updateManagementCMEAzureVwanProvision(d *schema.ResourceData, m interface{}) error {
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

	payload := make(map[string]interface{})

	if d.HasChange("base64_sic_key") {
		payload["base64_sic_key"] = d.Get("base64_sic_key").(string)
	}

	if d.HasChange("policy") {
		payload["policy"] = d.Get("policy").(string)
	}

	if d.HasChange("autonomous_threat_prevention"){
		payload["autonomous_threat_prevention"] = d.Get("autonomous_threat_prevention").(bool)
	}

	if d.HasChange("identity_awareness") {
		payload["identity_awareness"] = d.Get("identity_awareness").(bool)
	}

	log.Println("Update cme Azure VWAN provisioning - NVA = ", nvaName)

	url := CmeApiPath + "/azure/virtualWANs/accounts/" + accountID + "/resourceGroups/" + nvaResourceGroup + "/provision/" + nvaName

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

	return nil
}

func readManagementCMEAzureVwanProvision(d *schema.ResourceData, m interface{}) error {
	return nil
}

func deleteManagementCMEAzureVwanProvision(d *schema.ResourceData, m interface{}) error {
	d.SetId("")

	return nil
}
