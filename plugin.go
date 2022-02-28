package plugin

import (
	"github.com/corazawaf/coraza/v2"
	"github.com/corazawaf/coraza/v2/operators"
	"github.com/corazawaf/coraza/v2/seclang"
)

func init() {
	operators.RegisterPlugin("geoLookup", func() coraza.RuleOperator { return new(geoLookup) })
	seclang.RegisterDirectivePlugin("secGeoLookupDb", directiveSecGeoLookupDb)
}
