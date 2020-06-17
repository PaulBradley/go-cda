package cda

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

var dateFormat string
var dateTimeFormat string
var oidMappings = make(map[string]string)

// ClinicalDocument holds the parsed values
// from the supplied CDA XML document
type ClinicalDocument struct {
	htmlReport              string
	htmlReportLogoURL       string
	htmlReportStyleSheetURL string
	RawXML                  string `xml:",innerxml"`
	Title                   string `xml:"title"`

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

	StructuredBodySections []structuredBodySection `xml:"component>structuredBody>component"`
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
	v := 0
	for _, section := range cda.StructuredBodySections {
		cda.StructuredBodySections[v].AccreditedText = strings.Replace(section.AccreditedText, `\.br\`, replaceWith, -1)
		cda.StructuredBodySections[v].Text = strings.Replace(section.Text, `\.br\`, replaceWith, -1)
		v = v + 1
	}
}

// DoSetReportLogo sets a URL to a logo
// to be used within the HTML report.
func (cda *ClinicalDocument) DoSetReportLogo(URL string) {
	cda.htmlReportLogoURL = URL
}

// DoSetReportStyleSheet sets a URL to an external
// style sheet to be used within the HTML report.
func (cda *ClinicalDocument) DoSetReportStyleSheet(URL string) {
	cda.htmlReportStyleSheetURL = URL
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

	for _, section := range cda.StructuredBodySections {
		if len(section.EntryRelationshipOrganizer) > 0 {
			if len(section.EntryRelationshipOrganizer[0].EffectiveTime[0].Value) > 0 {
				dt, err = time.Parse("20060302150405-0700", section.EntryRelationshipOrganizer[0].EffectiveTime[0].Value)
				if err != nil {
					return err
				}
				section.EntryRelationshipOrganizer[0].EffectiveTime[0].Value = dt.Format(dateTimeFormat)
			}
		}
	}

	// date of birth
	dt, err = time.Parse("20060102", cda.RecordTarget[0].PatientRole.BirthTime[0].Value)
	if err != nil {
		return err
	}
	cda.RecordTarget[0].PatientRole.BirthTime[0].Value = dt.Format(dateFormat)

	return nil
}

// GenerateReport returns
// a formatted HTML 5 report
func (cda *ClinicalDocument) GenerateReport() string {

	cda.htmlReport = cda.webReportHeader()

	cda.webReportDocumentFields()
	cda.webReportPatient()

	cda.webReportDocumentSections()
	cda.webReportFooter()

	return cda.htmlReport
}

func (cda *ClinicalDocument) webReportHeader() string {
	return `<!doctype html>
	<html>
	<head>
		<meta charset="utf-8">
		<title>` + cda.Title + `</title>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="stylesheet" type="text/css" href="` + cda.htmlReportStyleSheetURL + `" integrity="sha384-9aIt2nRpC12Uk9gS9baDl411NQApFmC26EwAOH8WgZl5MYYxFfc+NcPb1dKGj7Sk" crossorigin="anonymous" />
		<style>
			body {
				margin: 2em 10%;
			}
			h3 {
				margin-top: 1em;
			}
			thead>tr>td {
				font-weight: bold;
			}
			caption {
				color: #444;
				font-weight: bold;
				caption-side: top;
			}
			@media print {
				main {
					display: none;
				}
			}
		</style>
	</head>
	<body>
	<main>
	<div class="container-fluid">`
}

func (cda *ClinicalDocument) webReportDocumentFields() {

	if len(cda.htmlReportLogoURL) > 0 {
		cda.htmlReport = cda.htmlReport + `<img src="` + cda.htmlReportLogoURL + `" style="float:right;" />`
	}

	cda.htmlReport = cda.htmlReport + `<h1>`
	if len(cda.Custodian[0].Name) > 0 {
		cda.htmlReport = cda.htmlReport + cda.Custodian[0].Name + `<br>`
	}
	if len(cda.Title) > 0 {
		cda.htmlReport = cda.htmlReport + cda.Title
	}
	cda.htmlReport = cda.htmlReport + `</h1>`

	if len(cda.EffectiveTime[0].Value) > 0 {
		cda.htmlReport = cda.htmlReport + `<p><b>Published ` + cda.EffectiveTime[0].Value + `</b></p>`
	}

	cda.htmlReport = cda.htmlReport + `<hr />`
}

func (cda *ClinicalDocument) webReportFooter() {
	cda.htmlReport = cda.htmlReport + `</div>
	</main>
	</body>
	</html>`
}

func (cda *ClinicalDocument) webReportPatient() {
	kcw := "15%"

	cda.htmlReport = cda.htmlReport + `<h3 class="text-primary">Patient</h3>`
	cda.htmlReport = cda.htmlReport + webReportTableOpen()
	cda.htmlReport = cda.htmlReport + webReportTableAddRow(kcw, "Patient", cda.DoFormatDisplayName(cda.RecordTarget[0].PatientRole.Person[0]))
	cda.htmlReport = cda.htmlReport + webReportTableAddRow(kcw, "Address", cda.DoFormatDisplayAddress(cda.RecordTarget[0].PatientRole.Address[0]))
	cda.htmlReport = cda.htmlReport + webReportTableAddRow(kcw, "DOB", cda.RecordTarget[0].PatientRole.BirthTime[0].Value)
	cda.htmlReport = cda.htmlReport + webReportTableAddRow(kcw, "Gender", cda.RecordTarget[0].PatientRole.AdministrativeGenderCode[0].DisplayName)
	cda.htmlReport = cda.htmlReport + webReportTableClose()

	kcw = "30%"
	cda.htmlReport = cda.htmlReport + `<h5 class="text-primary">Patient IDs</h5>`
	cda.htmlReport = cda.htmlReport + webReportTableOpen()

	for _, idNumbers := range cda.RecordTarget[0].PatientRole.IDs {
		cda.htmlReport = cda.htmlReport + webReportTableAddRow(kcw, idNumbers.Extension, idNumbers.Root)
	}

	cda.htmlReport = cda.htmlReport + webReportTableClose()
}

func (cda *ClinicalDocument) webReportDocumentSections() {
	for _, section := range cda.StructuredBodySections {
		cda.htmlReport = cda.htmlReport + `<h3 class="text-primary">` + section.Title + `</h3>`

		sectionHTML := section.Text
		sectionHTML = strings.Replace(sectionHTML, `<table>`, `<table class="table table-sm table-bordered">`, -1)

		cda.htmlReport = cda.htmlReport + sectionHTML

		if len(section.AccreditedText) > 0 {
			cda.htmlReport = cda.htmlReport + `<hr />`
			cda.htmlReport = cda.htmlReport + `<p>` + section.AccreditedText + `</p>`
			fmt.Println(section.AccreditedText)
		}
	}
}

func webReportTableOpen() string {
	return `<table id="patient" class="table table-sm">
	<tr>`
}

func webReportTableClose() string {
	return `</table>`
}

func webReportTableAddRow(width string, key string, value string) string {
	return `
	<tr>
		<td width="` + width + `">` + key + `</td>
		<td>` + value + `</td>
	</tr>`
}
