package kmgIpToCountry

import (
	"github.com/bronze1man/kmg/kmgTest"
	"net"
	"testing"
)

func TestGetCountryIsoCode(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	code, err := GetCountryIsoCode(net.ParseIP("180.97.33.107"))
	t.Equal(err, nil)
	t.Equal(code, "CN")

	code, err = GetCountryIsoCode(net.ParseIP("173.194.127.50"))
	t.Equal(err, nil)
	t.Equal(code, "US")
}
