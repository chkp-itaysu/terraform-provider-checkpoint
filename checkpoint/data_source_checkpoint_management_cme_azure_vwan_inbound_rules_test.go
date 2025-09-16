package checkpoint

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceCheckpointManagementCMEAzureInboundRules_basic(t *testing.T) {
	dataSourceName := "data.checkpoint_management_cme_azure_vwan_inbound_rules.test"
	accountId := "test-account"
	nvaResourceGroup := "test-nva-rg"
	nvaName := "test-nva"

	context := os.Getenv("CHECKPOINT_CONTEXT")
	if context == "" {
		t.Skip("Env CHECKPOINT_CONTEXT must be specified to run this test")
	} else if context != "web_api" {
		t.Skip("Skipping cme api test")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceManagementCMEAzureInboundRulesConfig(accountId, nvaResourceGroup, nvaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
				),
			},
		},
	})
}

func testAccDataSourceManagementCMEAzureInboundRulesConfig(accountId string, nvaResourceGroup string, nvaName string) string {
	return fmt.Sprintf(`
data "checkpoint_management_cme_azure_vwan_inbound_rules" "test" {
  account_id = "%s"
  nva_resource_group = "%s"
  nva_name = "%s"
}
`, accountId, nvaResourceGroup, nvaName)
}
