package cda

type addr struct {
	Country           string `xml:"country"`
	State             string `xml:"state"`
	City              string `xml:"city"`
	PostCode          string `xml:"postalCode"`
	StreetAddressLine string `xml:"streetAddressLine"`
}
