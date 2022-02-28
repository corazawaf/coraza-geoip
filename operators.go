package plugin

import (
	"net"
	"strconv"

	"github.com/corazawaf/coraza/v2"
	"github.com/corazawaf/coraza/v2/types/variables"
	geoip "github.com/oschwald/geoip2-golang"
)

type geoLookup struct{}

func (o *geoLookup) Init(data string) error {
	return nil
}

func (o *geoLookup) Evaluate(tx *coraza.Transaction, value string) bool {
	db, ok := tx.Waf.Config.Get("geoip", nil).(*geoip.Reader)
	if !ok || db == nil {
		return false
	}

	ip := net.ParseIP(value)
	r, err := db.City(ip)
	if err != nil {
		// do we generate an error?
		return false
	}

	tx.GetCollection(variables.Geo).Set("country_code", []string{r.Country.IsoCode})
	tx.GetCollection(variables.Geo).Set("country_name", []string{r.Country.Names["en"]})
	tx.GetCollection(variables.Geo).Set("country_continent", []string{r.Continent.Names["en"]})
	tx.GetCollection(variables.Geo).Set("region", []string{""})
	tx.GetCollection(variables.Geo).Set("city", []string{r.City.Names["en"]})
	tx.GetCollection(variables.Geo).Set("postal_code", []string{r.Postal.Code})
	tx.GetCollection(variables.Geo).Set("latitude", []string{strconv.FormatFloat(r.Location.Latitude, 'f', 10, 64)})
	tx.GetCollection(variables.Geo).Set("longitude", []string{strconv.FormatFloat(r.Location.Longitude, 'f', 10, 64)})
	return true
}

var _ coraza.RuleOperator = &geoLookup{}
