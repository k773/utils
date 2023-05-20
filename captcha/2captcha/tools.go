package twocaptcha

import (
	"github.com/k773/utils"
	"github.com/k773/utils/maps"
)

func structToMap(src any) map[string]string {
	return maps.ConvertValuesToString(utils.StructMapTagValueToFieldValue(src, "json", false, false))
}
