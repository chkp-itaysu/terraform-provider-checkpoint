---
layout: "checkpoint"
page_title: "checkpoint_management_cme_azure_vwan_provision"
sidebar_current: "docs-checkpoint-resource-checkpoint-management-cme-azure-vwan-provision"
description: |- This resource allows you to add/update/delete Check Point CME AWS Account.
---

# Resource: checkpoint_management_cme_azure_vwan_provision

This resource allows you to provision CME Azure VWAN NVA.

For details about the compatibility between the Terraform Release version and the CME API version, please refer to the section [Compatibility with CME](https://registry.terraform.io/providers/CheckPointSW/checkpoint/latest/docs#compatibility-with-cme).


## Example Usage

```hcl
resource "checkpoint_management_cme_azure_vwan_provision" "provision" {
  account_id                   = "azureAccount"
  nva_resource_group           = "nva_resource_group"
  nva_name                     = "nva_name"
  base64_sic_key               = "base64_sic_key"
  policy                       = "policy"
  autonomous_threat_prevention = true
  identity_awareness           = true
}
```

## Argument Reference

These arguments are supported:

* `account_id` - (Required) The ID of the Azure account.
* `nva_resource_group` - (Required) The name of the resource group that contains the NVA.
* `nva_name` - (Required) The name of the NVA.
* `base64_sic_key` - (Required) The sic key in base64 format.
* `policy` - (Required) The provision policy.
* `autonomous_threat_prevention` - (Required) Enable autonomous threat prevention.
* `identity_awareness` - (Required) Enable identity awareness.