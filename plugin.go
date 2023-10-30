package plugin

import (
	"github.com/corazawaf/coraza/v3/experimental/plugins"
	"github.com/oschwald/geoip2-golang"
)

func RegisterGeoDatabase(database []byte, databaseType string) error {
	return RegisterGeoDatabaseFromHandler(func() (*geoip2.Reader, error) {
		return geoip2.FromBytes(database)
	}, databaseType)
}

func RegisterGeoDatabaseFromFile(databasePath string, databaseType string) error {
	return RegisterGeoDatabaseFromHandler(func() (*geoip2.Reader, error) {
		return geoip2.Open(databasePath)
	}, databaseType)
}

func RegisterGeoDatabaseFromHandler(handler func() (*geoip2.Reader, error), databaseType string) error {
	db, err := handler()
	if err != nil {
		return err
	}
	plugins.RegisterOperator("geoLookup", newGeolookupBuilder(db, databaseType))
	return nil
}
