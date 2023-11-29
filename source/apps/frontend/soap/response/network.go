package response

import "encoding/xml"

type Envelope2 struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Header  struct {
		Text           string `xml:",chardata"`
		ResponseHeader struct {
			Text         string `xml:",chardata"`
			Xmlns        string `xml:"xmlns,attr"`
			RequestId    string `xml:"requestId"`
			ResponseTime string `xml:"responseTime"`
		} `xml:"ResponseHeader"`
	} `xml:"Header"`
	Body struct {
		Text                   string `xml:",chardata"`
		GetAllNetworksResponse struct {
			Text  string `xml:",chardata"`
			Xmlns string `xml:"xmlns,attr"`
			Rval  []struct {
				Text                  string `xml:",chardata"`
				ID                    string `xml:"id"`
				DisplayName           string `xml:"displayName"`
				NetworkCode           string `xml:"networkCode"`
				PropertyCode          string `xml:"propertyCode"`
				TimeZone              string `xml:"timeZone"`
				CurrencyCode          string `xml:"currencyCode"`
				EffectiveRootAdUnitId string `xml:"effectiveRootAdUnitId"`
				IsTest                string `xml:"isTest"`
			} `xml:"rval"`
		} `xml:"getAllNetworksResponse"`
	} `xml:"Body"`
}


type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Soap    string   `xml:"soap,attr"`
	Header  struct {
		Text           string `xml:",chardata"`
		ResponseHeader struct {
			Text         string `xml:",chardata"`
			Xmlns        string `xml:"xmlns,attr"`
			RequestId    string `xml:"requestId"`
			ResponseTime string `xml:"responseTime"`
		} `xml:"ResponseHeader"`
	} `xml:"Header"`
	Body struct {
		Text                   string `xml:",chardata"`
		GetAllNetworksResponse struct {
			Text  string `xml:",chardata"`
			Xmlns string `xml:"xmlns,attr"`
			Rval  struct {
				Text                  string `xml:",chardata"`
				ID                    string `xml:"id"`
				DisplayName           string `xml:"displayName"`
				NetworkCode           string `xml:"networkCode"`
				PropertyCode          string `xml:"propertyCode"`
				TimeZone              string `xml:"timeZone"`
				CurrencyCode          string `xml:"currencyCode"`
				EffectiveRootAdUnitId string `xml:"effectiveRootAdUnitId"`
				IsTest                string `xml:"isTest"`
			} `xml:"rval"`
		} `xml:"getAllNetworksResponse"`
	} `xml:"Body"`
}
