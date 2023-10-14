# GeoIP Lookup Plugin for Coraza WAF

## Introduction

This plugin enables GeoIP lookup functionality for Coraza Web Application Firewall (WAF) by utilizing a GeoIP database. With this, you can efficiently track and filter requests based on geographical origins.

## Features

- Seamless integration with Coraza WAF.
- Efficient IP to geographical location mapping.
- Customizable actions based on country, region, or city.

## Prerequisites

- Coraza WAF installed and configured.
- GeoIP Database (e.g., MaxMind GeoLite2 or any compatible database).

## Installation

1. Add the go module to your project:

   ```bash
   go get https://github.com/corazawaf/coraza-geoip
   ```

2. Import the plugin in your project and configure it (supported database types: `country`, `city`):

   #### Using an embedded database file (if local file access is not available, e.g. WebAssembly)

   ```go
   import (
        _ "embed"
       geo "github.com/corazawaf/coraza-geoip"
   )

   //go:embed geoip-database.mmdb
   var geoIpDatabase []byte

   func init() {
       geo.RegisterDatabase(geoIpDatabase, "country")
   }
   ```

   #### Using a local database file

   ```go
      import (
        _ "embed"
       geo "github.com/corazawaf/coraza-geoip"
   )


   func init() {
       geo.RegisterDatabaseFromFile("geoip-database.mmdb", "city")
   }
   ```

## Usage

After setting up, Coraza WAF will start utilizing the GeoIP database to determine the geographical location of incoming IP addresses. You can then create rules to allow, block, or log requests based on their geographic origin.

See also the [Coraza WAF documentation about using "@geoLookup"](https://coraza.io/docs/seclang/operators/#geolookup) and [Using Plugins documentation](https://coraza.io/docs/tutorials/using-plugins/) for more information.
