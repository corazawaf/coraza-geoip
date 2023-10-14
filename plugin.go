package plugin

import (
	"github.com/corazawaf/coraza/v3/experimental/plugins"
	"github.com/oschwald/geoip2-golang"
)

func RegisterGeoDatabase(database []byte, databaseType string) {
	db, err := geoip2.FromBytes(database)
	if err != nil {
		panic(err)
	}

	plugins.RegisterOperator("geoLookup", newGeolookupCreator(db, databaseType))
}

func RegisterGeoDatabaseFromFile(databasePath string, databaseType string) {
	db, err := geoip2.Open(databasePath)
	if err != nil {
		panic(err)
	}

	plugins.RegisterOperator("geoLookup", newGeolookupCreator(db, databaseType))
}
