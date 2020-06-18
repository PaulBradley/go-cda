package cda

type codeSystem struct {
	Code           string `xml:"code,attr"`
	CodeSystem     string `xml:"codeSystem,attr"`
	CodeSystemName string `xml:"codeSystemName,attr"`
	DisplayName    string `xml:"displayName,attr"`
}
