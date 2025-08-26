---
layout: "checkpoint"
page_title: "checkpoint_management_cme_azure_vwan_inbound_rules"
sidebar_current: "docs-checkpoint-data-source-checkpoint-management-cme-azure-vwan-inbound-rules"
description: |- Use this data source to get information on all Check Point CME Azure VWAN Inbound Rules.
---

# Data Source: checkpoint_management_cme_azure_vwan_inbound_rules

Use this data source to get information on all Check Point CME Azure VWAN Inbound Rules.

For details about the compatibility between the Terraform Release version and the CME API version, please refer to the section [Compatibility with CME](https://registry.terraform.io/providers/CheckPointSW/checkpoint/latest/docs#compatibility-with-cme).


## Example Usage

```hcl
data "checkpoint_management_cme_azure_vwan_inbound_rules" "incound_rules" {
    account_id         = "CME Azure Account ID"
    nva_resource_group = "Azure NVA resource group name"
    nva_name           = "Azure NVA resource name"
}
```

## Argument Reference

These arguments are supported:

* `account_id` - (Required) The ID of the Azure account.
* `nva_resource_group` - (Required) The name of the resource group that contains the NVA.
* `nva_name` - (Required) The name of the NVA.
* `rules` - List of all inbound rules, each with this data:
    * `name` - Unique rule name for identification.
    * `lb_public_ips` - List of outbound public IPs for that rule.
    * `original_ports` - The list of ports allowed in that rule.
    * `original_source` - The rule inbound IP address.
    * `protocol` - The traffic protocol in that rule.