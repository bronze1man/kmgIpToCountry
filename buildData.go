package kmgIpToCountry
//this file is generate by script, do not modify it by hand.
import "encoding/base64"
func getData()[]byte{
	gzContent, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			panic(err)
		}
		return gzContent
}