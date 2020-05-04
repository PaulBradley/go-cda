package cda

type structuredBodySection struct {
	Title   string `xml:"section>title"`
	Text    string `xml:"section>text"`
	ActText string `xml:"section>entry>organizer>component>act>text"`
}
