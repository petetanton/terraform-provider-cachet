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

- `description` (String)
- `name` (String)

### Optional

- `enabled` (Boolean)
- `group_id` (Number)
- `link` (String)
- `status` (String)
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `default` (String)
- `delete` (String)
- `read` (String)
- `update` (String)

