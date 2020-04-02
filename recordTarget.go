package cda

type administrativeGenderCode struct {
	Code           string `xml:"code,attr"`
	CodeSystem     string `xml:"codeSystem,attr"`
	CodeSystemName string `xml:"codeSystemName,attr"`
	DisplayName    string `xml:"displayName,attr"`
}

type birthTime struct {
	Value string `xml:"value,attr"`
}

type patientRole struct {
	IDs                      []id                       `xml:"id"`
	Address                  []addr                     `xml:"addr"`
	Person                   []person                   `xml:"patient>name"`
	BirthTime                []birthTime                `xml:"patient>birthTime"`
	AdministrativeGenderCode []administrativeGenderCode `xml:"patient>administrativeGenderCode"`
	MaritalStatusCode        string                     `xml:"patient>maritalStatusCode"`
	ReligiousAffiliationCode string                     `xml:"patient>religiousAffiliationCode"`
}

type recordTarget struct {
	TypeCode    string      `xml:"typeCode,attr"`
	PatientRole patientRole `xml:"patientRole"`
}
