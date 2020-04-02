package cda

type custodian struct {
	Name    string `xml:"name"`
	Address []addr `xml:"addr"`
}
