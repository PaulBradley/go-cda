package cda

import (
	"encoding/xml"
	"strings"
	"time"
)

var dateFormat string
var dateTimeFormat string

// ClinicalDocument holds the parsed values
// from the supplied CDA XML document
type ClinicalDocument struct {
	RawXML string `xml:",innerxml"`
	Title  string `xml:"title"`

	// document header items with attributes
	Code          []codeSystem    `xml:"code"`
	EffectiveTime []effectiveTime `xml:"effectiveTime"`
	IDs           []id            `xml:"id"`
	LanguageCode  []languageCode  `xml:"languageCode"`
	RealmCode     []realmCode     `xml:"realmCode"`
	VersionNumber []versionNumber `xml:"versionNumber"`

	// patient
	RecordTarget []recordTarget `xml:"recordTarget"`

	// document custodian and author(s)
	Custodian []custodian `xml:"custodian>assignedCustodian>representedCustodianOrganization"`
	Authors   []author    `xml:"author"`

	// recipient(s)
	IntendedRecipients []intendedRecipient `xml:"informationRecipient>intendedRecipient"`

	// encounter information
	EncompassingEncounter []encompassingEncounter `xml:"componentOf>encompassingEncounter"`

	// unstructured content
	StructuredBodyTitle          string `xml:"component>structuredBody>component>section>title"`
	StructuredBodyText           string `xml:"component>structuredBody>component>section>text"`
	StructuredBodyTextAdditional string `xml:"component>structuredBody>component>section>entry>organizer>component>act>text"`

	// lab specific
	LabAccessionNumber []id            `xml:"component>structuredBody>component>section>entry>act>entryRelationship>organizer>id"`
	LabStatusCode      []code          `xml:"component>structuredBody>component>section>entry>act>entryRelationship>organizer>statusCode"`
	LabEffectiveTime   []effectiveTime `xml:"component>structuredBody>component>section>entry>act>entryRelationship>organizer>effectiveTime"`
}

// Parse accepts the CDA document as a string
// and passes it to the Unmarshal function to
// map the data items into in the struct fields.
func Parse(clinicalDocument string) (ClinicalDocument, error) {

	parsedClinicDocument := ClinicalDocument{}

	// wrap cdata tag around embedded <text> nodes
	clinicalDocument = strings.Replace(clinicalDocument, "<text>", "<text><![CDATA[", -1)
	clinicalDocument = strings.Replace(clinicalDocument, "</text>", "]]></text>", -1)

	err := xml.Unmarshal([]byte(clinicalDocument), &parsedClinicDocument)
	if err != nil {
		return parsedClinicDocument, err
	}

	return parsedClinicDocument, nil
}

// DoSwitchBR replaces the carriage return text
// within the LabReports from \.br\ to the string
// supplied in the replaceWith parameter
func (cda *ClinicalDocument) DoSwitchBR(replaceWith string) {
	cda.StructuredBodyText = strings.Replace(cda.StructuredBodyText, `\.br\`, replaceWith, -1)
	cda.StructuredBodyTextAdditional = strings.Replace(cda.StructuredBodyTextAdditional, `\.br\`, replaceWith, -1)
}

// DoFormatDisplayName formats a name by
// concatenating Prefix, Given, Family
// and Suffix values
func (cda *ClinicalDocument) DoFormatDisplayName(p person) string {
	var displayName string

	displayName = ""
	if len(p.Prefix) > 0 {
		displayName = displayName + p.Prefix
	}
	if len(p.Given) > 0 {
		displayName = displayName + " " + p.Given
	}
	if len(p.Family) > 0 {
		displayName = displayName + " " + p.Family
	}
	if len(p.Suffix) > 0 {
		displayName = displayName + " " + p.Suffix
	}

	return displayName
}

// DoFormatDisplayAddress formats an address
// by concatenating the address segments
func (cda *ClinicalDocument) DoFormatDisplayAddress(a addr) string {
	var displayAddress string

	displayAddress = ""
	if len(a.StreetAddressLine) > 0 && a.StreetAddressLine != "NA" {
		displayAddress = displayAddress + a.StreetAddressLine + "<br />"
	}
	if len(a.City) > 0 && a.City != "NA" {
		displayAddress = displayAddress + a.City + "<br />"
	}
	if len(a.Country) > 0 && a.Country != "NA" {
		displayAddress = displayAddress + a.Country + "<br />"
	}
	if len(a.State) > 0 && a.State != "NA" {
		displayAddress = displayAddress + a.State + "<br />"
	}
	if len(a.PostCode) > 0 && a.PostCode != "NA" {
		displayAddress = displayAddress + a.PostCode
	}

	return displayAddress
}

// DoSetDateFormat sets the format to be
// applied to date fields like DOB
func (cda *ClinicalDocument) DoSetDateFormat(format string) {
	dateFormat = format
}

// DoSetDateTimeFormat sets the format
// to be applied to date/time fields like
// effectiveTime
func (cda *ClinicalDocument) DoSetDateTimeFormat(format string) {
	dateTimeFormat = format
}

// DoReformatDateTimeFields parses the date/time
// fields and applies the formatting mask passed
// via DoSetDateFormat() & DoSetDateTimeFormat()
func (cda *ClinicalDocument) DoReformatDateTimeFields() error {
	var err error
	var dt time.Time

	// document effectiveTime
	dt, err = time.Parse("20060302150405-0700", cda.EffectiveTime[0].Value)
	if err != nil {
		return err
	}
	cda.EffectiveTime[0].Value = dt.Format(dateTimeFormat)

	// lab effectiveTime
	dt, err = time.Parse("20060302150405-0700", cda.LabEffectiveTime[0].Value)
	if err != nil {
		return err
	}
	cda.LabEffectiveTime[0].Value = dt.Format(dateTimeFormat)

	// date of birth
	dt, err = time.Parse("20060102", cda.RecordTarget[0].PatientRole.BirthTime[0].Value)
	if err != nil {
		return err
	}
	cda.RecordTarget[0].PatientRole.BirthTime[0].Value = dt.Format(dateFormat)

	return nil
}
