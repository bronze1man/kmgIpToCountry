package kmgIpToCountry

import (
	"bytes"
	"compress/gzip"
	"github.com/oschwald/geoip2-golang"
	"io/ioutil"
	"net"
	"sync"
)

func Country(ip net.IP) (*geoip2.Country, error) {
	EnsureInit()
	return geoip2Reader.Country(ip)
}

// iso code look like CN, US
func GetCountryIsoCode(ip net.IP) (code string, err error) {
	EnsureInit()
	c, err := geoip2Reader.Country(ip)
	if err != nil {
		return
	}
	return c.Country.IsoCode, nil
}

// iso code look like CN, US
func MustGetCountryIsoCode(ip net.IP) (code string) {
	EnsureInit()
	c, err := geoip2Reader.Country(ip)
	if err != nil {
		panic(err)
	}
	return c.Country.IsoCode
}

var geoip2Reader *geoip2.Reader
var geoipInitOnce sync.Once

func EnsureInit() {
	geoipInitOnce.Do(func() {
		gzContent := getData()
		r := bytes.NewReader(gzContent)
		gzReader, err := gzip.NewReader(r)
		if err != nil {
			panic(err)
		}
		uncompressData, err := ioutil.ReadAll(gzReader)
		if err != nil {
			panic(err)
		}
		geoip2Reader, err = geoip2.FromBytes(uncompressData)
		if err != nil {
			panic(err)
		}
		gzContent = nil
	})
}
