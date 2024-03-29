---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cachet_component Resource - terraform-provider-cachet"
subcategory: ""
description: |-
  A component is a resource that defines a deployable thing that can be operational or degraded
---

# cachet_component (Resource)

A component is a resource that defines a deployable thing that can be operational or degraded



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `description` (String) Description of the component
- `name` (String) Name of the component

### Optional

- `enabled` (Boolean) Is the component enabled
- `group_id` (Number) The group that the component is within
- `link` (String) A link to the component - can be used for Runbooks etc
- `status` (String) The current status of thc component. If you have automation that updates the component status, please use lifecycle rules in terraform

### Read-Only

- `id` (String) The ID of this resource.


