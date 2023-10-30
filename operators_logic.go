package plugin

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/corazawaf/coraza/v3/experimental/plugins/plugintypes"
	"github.com/oschwald/geoip2-golang"
)

type geo struct {
	db     geoIPReader
	dbtype string
}

func newGeolookupBuilder(db *geoip2.Reader, databaseType string) func(options plugintypes.OperatorOptions) (plugintypes.Operator, error) {
	return func(options plugintypes.OperatorOptions) (plugintypes.Operator, error) {
		return newGeolookup(options, db, databaseType)
	}
}

func newGeolookup(options plugintypes.OperatorOptions, db *geoip2.Reader, databaseType string) (plugintypes.Operator, error) {
	return &geo{db: db, dbtype: databaseType}, nil
}

func (o *geo) applyVariablesCity(col mapCollection, ip net.IP) (bool, error) {
	r, err := o.db.City(ip)
	if err != nil {
		return false, err
	}

	col.Set("country_code", []string{r.Country.IsoCode})
	col.Set("country_name", []string{r.Country.Names["en"]})
	col.Set("continent_code", []string{r.Continent.Code})
	col.Set("country_continent", []string{r.Continent.Names["en"]})
	col.Set("region", []string{""})
	col.Set("city", []string{r.City.Names["en"]})
	col.Set("postal_code", []string{r.Postal.Code})
	col.Set("latitude", []string{strconv.FormatFloat(r.Location.Latitude, 'f', 10, 64)})
	col.Set("longitude", []string{strconv.FormatFloat(r.Location.Longitude, 'f', 10, 64)})

	return true, nil
}

func (o *geo) applyVariablesCountry(col mapCollection, ip net.IP) (bool, error) {
	r, err := o.db.Country(ip)
	if err != nil {
		return false, err
	}

	col.Set("country_code", []string{r.Country.IsoCode})
	col.Set("country_name", []string{r.Country.Names["en"]})
	col.Set("continent_code", []string{r.Continent.Code})
	col.Set("country_continent", []string{r.Continent.Names["en"]})

	return true, nil
}

func (o *geo) executeEvaluationInternal(tx transaction, value string) (bool, error) {
	if len(value) == 0 {
		return false, errors.New("empty IP")
	}
	col, err := tx.GetGeoCollection(tx)
	if err != nil {
		return false, err
	}
	ip := net.ParseIP(value)
	if ip == nil {
		return false, fmt.Errorf("invalid ip %q", value)
	}

	if ip.IsPrivate() || ip.IsLoopback() || ip.IsUnspecified() {
		return false, nil
	}

	switch o.dbtype {
	case "city":
		return o.applyVariablesCity(col, ip)
	case "country":
		return o.applyVariablesCountry(col, ip)
	default:
		return false, fmt.Errorf("invalid database type")
	}
}
