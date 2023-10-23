package plugin

import (
	"net"
	"testing"

	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApplyVariablesCity(t *testing.T) {
	sut := setupTest("city")

	ip := net.ParseIP("8.8.8.8")
	result, err := sut.geoInstance.applyVariablesCity(sut.col, ip)

	assert.NoError(t, err)
	assert.True(t, result)

	sut.col.AssertCalled(t, "Set", "country_code", []string{"US"})
	sut.col.AssertCalled(t, "Set", "country_name", []string{"United States"})
	sut.col.AssertCalled(t, "Set", "country_continent", []string{"North America"})
	sut.col.AssertCalled(t, "Set", "region", []string{""})
	sut.col.AssertCalled(t, "Set", "city", []string{"New York"})
	sut.col.AssertCalled(t, "Set", "postal_code", []string{"10001"})
	sut.col.AssertCalled(t, "Set", "latitude", []string{"40.7306100000"})
	sut.col.AssertCalled(t, "Set", "longitude", []string{"-73.9352420000"})
}

func TestApplyVariablesCountry(t *testing.T) {
	sut := setupTest("country")

	ip := net.ParseIP("8.8.8.8")
	result, err := sut.geoInstance.applyVariablesCountry(sut.col, ip)

	assert.NoError(t, err)
	assert.True(t, result)

	sut.col.AssertCalled(t, "Set", "country_code", []string{"US"})
	sut.col.AssertCalled(t, "Set", "country_name", []string{"United States"})
	sut.col.AssertCalled(t, "Set", "country_continent", []string{"North America"})
}

func TestEvaluateCity(t *testing.T) {
	sut := setupTest("city")

	_, err := sut.geoInstance.executeEvaluationInternal(sut.tx, "8.8.8.8")

	assert.NoError(t, err)
	sut.reader.AssertCalled(t, "City", net.IPv4(8, 8, 8, 8))
}

func TestEvaluateCountry(t *testing.T) {
	sut := setupTest("country")

	_, err := sut.geoInstance.executeEvaluationInternal(sut.tx, "8.8.8.8")

	assert.NoError(t, err)
	sut.reader.AssertCalled(t, "Country", net.IPv4(8, 8, 8, 8))
}

func TestEvaluateInvalidDbType(t *testing.T) {
	sut := setupTest("invalid")

	_, err := sut.geoInstance.executeEvaluationInternal(sut.tx, "8.8.8.8")

	assert.Error(t, err)
}

func TestReturnFalseOnInvalidIp(t *testing.T) {
	sut := setupTest("city")

	_, err := sut.geoInstance.executeEvaluationInternal(sut.tx, "invalid")

	assert.Error(t, err)
}

type mockGeoIPReader struct {
	mock.Mock
}

func newMockGeoIPReader() *mockGeoIPReader {
	reader := &mockGeoIPReader{}
	reader.On("City", mock.Anything).Return(&geoip2.City{
		Country: struct {
			Names             map[string]string "maxminddb:\"names\""
			IsoCode           string            "maxminddb:\"iso_code\""
			GeoNameID         uint              "maxminddb:\"geoname_id\""
			IsInEuropeanUnion bool              "maxminddb:\"is_in_european_union\""
		}{
			IsoCode: "US",
			Names:   map[string]string{"en": "United States"},
		},
		Continent: struct {
			Names     map[string]string "maxminddb:\"names\""
			Code      string            "maxminddb:\"code\""
			GeoNameID uint              "maxminddb:\"geoname_id\""
		}{
			Code:  "NA",
			Names: map[string]string{"en": "North America"},
		},
		City: struct {
			Names     map[string]string "maxminddb:\"names\""
			GeoNameID uint              "maxminddb:\"geoname_id\""
		}{
			Names: map[string]string{"en": "New York"},
		},
		Postal: struct {
			Code string "maxminddb:\"code\""
		}{
			Code: "10001",
		},
		Location: struct {
			TimeZone       string  "maxminddb:\"time_zone\""
			Latitude       float64 "maxminddb:\"latitude\""
			Longitude      float64 "maxminddb:\"longitude\""
			MetroCode      uint    "maxminddb:\"metro_code\""
			AccuracyRadius uint16  "maxminddb:\"accuracy_radius\""
		}{
			Latitude:  40.730610,
			Longitude: -73.935242,
		},
	}, nil)

	reader.On("Country", mock.Anything).Return(&geoip2.Country{
		Country: struct {
			Names             map[string]string "maxminddb:\"names\""
			IsoCode           string            "maxminddb:\"iso_code\""
			GeoNameID         uint              "maxminddb:\"geoname_id\""
			IsInEuropeanUnion bool              "maxminddb:\"is_in_european_union\""
		}{
			IsoCode: "US",
			Names:   map[string]string{"en": "United States"},
		},
		Continent: struct {
			Names     map[string]string "maxminddb:\"names\""
			Code      string            "maxminddb:\"code\""
			GeoNameID uint              "maxminddb:\"geoname_id\""
		}{
			Code:  "NA",
			Names: map[string]string{"en": "North America"},
		},
	}, nil)

	return reader
}

func (r *mockGeoIPReader) City(ip net.IP) (*geoip2.City, error) {
	args := r.Called(ip)
	return args.Get(0).(*geoip2.City), args.Error(1)
}

func (r *mockGeoIPReader) Country(ip net.IP) (*geoip2.Country, error) {
	args := r.Called(ip)
	return args.Get(0).(*geoip2.Country), args.Error(1)
}

type collectionMock struct {
	mock.Mock
}

func newCollectionMock() *collectionMock {
	colMock := &collectionMock{}
	colMock.On("Set", mock.Anything, mock.Anything).Return()
	return colMock
}

func (c *collectionMock) Set(key string, values []string) {
	c.Called(key, values)
}

func setupTest(dbType string) *sut {
	reader := newMockGeoIPReader()
	geoInstance := &geo{db: reader, dbtype: dbType}
	col := newCollectionMock()
	tx := newTxMock(col)
	return &sut{geoInstance: geoInstance, tx: tx, reader: reader, col: col}
}

type transactionMock struct {
	mock.Mock
}

func newTxMock(col *collectionMock) *transactionMock {
	txMock := &transactionMock{}
	txMock.On("GetGeoCollection", mock.Anything).Return(col, nil)
	return txMock
}

func (t *transactionMock) GetGeoCollection(tx transaction) (mapCollection, error) {
	args := t.Called(tx)
	return args.Get(0).(mapCollection), args.Error(1)
}

type sut struct {
	geoInstance *geo
	reader      *mockGeoIPReader
	tx          *transactionMock
	col         *collectionMock
}
