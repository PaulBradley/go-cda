package cda

type reportedTime struct {
	Value string `xml:"value,attr"`
}

type representedOrganization struct {
	IDs      []id      `xml:"id"`
	Name     string    `xml:"name"`
	Addr     []addr    `xml:"addr"`
	Telecoms []telecom `xml:"telecom"`
}

type author struct {
	ReportedTime            []reportedTime            `xml:"time"`
	IDs                     []id                      `xml:"assignedAuthor>id"`
	Telecoms                []telecom                 `xml:"assignedAuthor>telecom"`
	Address                 []addr                    `xml:"assignedAuthor>addr"`
	AssignedPerson          []person                  `xml:"assignedAuthor>assignedPerson>name"`
	RepresentedOrganization []representedOrganization `xml:"assignedAuthor>representedOrganization"`
	ManufacturerModelName   string                    `xml:"assignedAuthor>assignedAuthoringDevice>manufacturerModelName"`
	SoftwareName            string                    `xml:"assignedAuthor>assignedAuthoringDevice>softwareName"`
}
