package plugin

import (
	"fmt"
	"net"

	"github.com/corazawaf/coraza/v3/collection"
	"github.com/corazawaf/coraza/v3/experimental/plugins/plugintypes"
	"github.com/corazawaf/coraza/v3/types/variables"
	"github.com/oschwald/geoip2-golang"
)

var _ transaction = (*transactionWrapper)(nil)

// this acts as a shim for the geoip2.Reader struct so we can mock it in tests
type geoIPReader interface {
	City(ip net.IP) (*geoip2.City, error)
	Country(ip net.IP) (*geoip2.Country, error)
}

// this acts as a shim for the collection.Map struct so we can mock it in tests
type mapCollection interface {
	Set(key string, values []string)
}

// this acts as a shim for the plugintypes.TransactionState struct so we can mock it in tests
type transaction interface {
	GetGeoCollection(tx transaction) (mapCollection, error)
}

type transactionWrapper struct {
	tx plugintypes.TransactionState
}

func (t *transactionWrapper) GetGeoCollection(tx transaction) (mapCollection, error) {
	if c, ok := t.tx.Collection(variables.Geo).(collection.Map); ok {
		if c == nil {
			return nil, fmt.Errorf("collection is nil")
		}
		return c, nil
	}
	return nil, fmt.Errorf("collection not found or not a map")
}

// this is the entry point for the plugin, we hide the implementation details here
// and expose a slimmer interface to the actual plugin logic
func (o *geo) Evaluate(tx plugintypes.TransactionState, value string) bool {
	transaction := &transactionWrapper{tx: tx}
	result, err := o.executeEvaluationInternal(transaction, value)
	if err != nil {
		tx.DebugLogger().Error().Msg(fmt.Sprintf("error looking up geoip: %s", err))
		return false
	}
	return result
}
