package cda

type effectiveTime struct {
	Value string `xml:"value,attr"`
}

type languageCode struct {
	Code string `xml:"code,attr"`
}

type realmCode struct {
	Code string `xml:"code,attr"`
}

type versionNumber struct {
	VersionNumber string `xml:"value,attr"`
}
