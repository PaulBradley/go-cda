package cda

type receivedOrganization struct {
	Address  []addr    `xml:"addr"`
	Name     string    `xml:"name"`
	Telecoms []telecom `xml:"telecom"`
}

type intendedRecipient struct {
	Address               []addr                 `xml:"addr"`
	InformationRecipient  []person               `xml:"informationRecipient>name"`
	ReceivedOrganizations []receivedOrganization `xml:"receivedOrganization"`
	Telecoms              []telecom              `xml:"telecom"`
}
