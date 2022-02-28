package plugin

import (
	"github.com/corazawaf/coraza/v2"
	"github.com/oschwald/geoip2-golang"
)

func directiveSecGeoLookupDb(w *coraza.Waf, opts string) error {
	db, err := geoip2.Open(opts)
	if err != nil {
		return err
	}
	w.Config.Set("geoip", db)
	return nil
}

// directive type was not exported
// var _ coraza.Directive = directiveSecGeoLookupDb
