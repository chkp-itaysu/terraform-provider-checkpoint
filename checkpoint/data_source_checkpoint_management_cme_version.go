package checkpoint

import (
	"fmt"

	checkpoint "github.com/CheckPointSW/cp-mgmt-api-go-sdk/APIFiles"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceManagementCMEVersion() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceManagementCMEVersionRead,
		Schema: map[string]*schema.Schema{},
	}
}

func dataSourceManagementCMEVersionRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*checkpoint.ApiClient)
	payload := map[string]interface{}{}
	cmeVersionRes, err := client.ApiCall("cme-api/v1/generalConfiguration/cmeVersion", payload, client.GetSessionID(), true, client.IsProxyUsed())

	if err != nil {
		return fmt.Errorf(err.Error())
	}
	if !cmeVersionRes.Success {
		return fmt.Errorf(cmeVersionRes.ErrorMsg)
	}
	ruleBaseJson := cmeVersionRes.GetData()
	fmt.Println(ruleBaseJson)

	return nil
}
