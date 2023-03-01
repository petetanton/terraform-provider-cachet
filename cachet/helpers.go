package cachet

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func StringInMapKeys(in map[string]interface{}, ignoreCase bool) schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.StringInSlice(GetMapKeys(in), ignoreCase))
}

func FindIntInMap(m map[string]interface{}, i int, def string) string {
	for name, j := range m {
		if i == j.(int) {
			return name
		}
	}

	return def
}

func GetMapKeys(in map[string]interface{}) []string {
	var out []string

	for s, _ := range in {
		out = append(out, s)
	}

	return out
}
