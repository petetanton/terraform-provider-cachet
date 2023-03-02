---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cachet_metric Resource - terraform-provider-cachet"
subcategory: ""
description: |-
  A component group is a resource that defines a group of components
---

# cachet_metric (Resource)

A component group is a resource that defines a group of components



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `decimal_places` (Number)
- `default_value` (Number)
- `default_view` (String)
- `description` (String)
- `mins_between_datapoints` (Number)
- `name` (String)
- `unit` (String)
- `visibility` (String)

### Optional

- `calculation_type` (String)
- `display_chart` (Boolean)
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

