---
layout: "checkpoint"
page_title: "checkpoint_management_cme_azure_vwan_inbound_rules"
sidebar_current: "docs-checkpoint-resource-checkpoint-management-cme-azure-vwan-inbound-rules"
description: |-
  This resource allows you to provision Azure VWAN NVA.
---

# Resource: checkpoint_management_cme_azure_vwan_inbound_rules

This resource allows you to add/delete Azure VWAN NVA inbound rules.

For details about the compatibility between the Terraform Release version and the CME API version, please refer to the section [Compatibility with CME](https://registry.terraform.io/providers/CheckPointSW/checkpoint/latest/docs#compatibility-with-cme).


## Example Usage

```hcl
resource "checkpoint_management_cme_azure_vwan_inbound_rules" "rules" {
  account_id         = "azureAccount"
  nva_resource_group = "nva_resource_group"
  nva_name           = "nva_name"
  rules              = "[{\"name\": \"exmaple\", \"original_source\": \"0.0.0.0/0\", \"lb_public_ips\": [\"123.123.123.123\"], \"original_ports\": [\"80\"], \"protocol\": \"TCP\"}]"
}
```

*Note:* For more informations about the rules parameter checkout the [Rules](#rules) section.

## Rules
The rules parameter is a string in a list of maps format, for example a list of rules will look like this:
```
...
rules = "[{\"name\": \"exmaple\", \"original_source\": \"0.0.0.0/0\", \"lb_public_ips\": [\"123.123.123.123\"], \"original_ports\": [\"80\"], \"protocol\": \"TCP\"}]"
...
```

Rule is a map with the following structure:
```
{
  "name"            : string,
  "original_source" : string,
  "lb_public_ips"   : list(string),
  "original_ports"  : list(string),
  "protocol"        : string
}
```

For example:
```json
{
  "name"            : "exmaple",
  "original_source" : "0.0.0.0/0",
  "lb_public_ips"   : ["123.123.123.123"],
  "original_ports"  : ["80"],
  "protocol"        : "TCP"
}
```

## Best Practice
When you want to define rules a good practice is to build a list of rules and to encode it using the terraform built in function `jsonencode`:
```hcl
rules = jsonencode([
    {
        "name"            : "example-1",
        "original_source" : "0.0.0.0/0",
        "lb_public_ips"   : ["123.123.123.123"],
        "original_ports"  : ["80"],
        "protocol"        : "TCP"
    },
    {
        "name"            : "example-2",
        "original_source" : "0.0.0.0/0",
        "lb_public_ips"   : ["123.123.123.123"],
        "original_ports"  : ["443"],
        "protocol"        : "TCP"
    },
])
```

The same goes for the output, when retrieving the output of the rules you can use the terraform built in function `jsondecode` to parse the string:
```hcl
output "rules" {
  value = jsondecode(resource.checkpoint_management_cme_azure_vwan_inbound_rules.example.rules)
}
```

## Argument Reference
These arguments are supported:
* `account_id` - (Required) The ID of the Azure account.
* `nva_resource_group` - (Required) The name of the resource group that contains the NVA.
* `nva_name` - (Required) The name of the NVA.
* `rules` - List of all inbound rules in a string format, each with this data:
    * `name` - Unique rule name for identification.
    * `lb_public_ips` - List of outbound public IPs for that rule.
    * `original_ports` - The list of ports allowed in that rule.
    * `original_source` - The rule inbound IP address.
    * `protocol` - The traffic protocol in that rule.
