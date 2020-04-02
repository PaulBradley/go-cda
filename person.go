package cda

type person struct {
	Given  string `xml:"given"`
	Family string `xml:"family"`
	Suffix string `xml:"suffix"`
	Prefix string `xml:"prefix"`
}
