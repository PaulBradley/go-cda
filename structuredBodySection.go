package cda

type structuredBodySection struct {
	Title                      string      `xml:"section>title"`
	Text                       string      `xml:"section>text"`
	AccreditedText             string      `xml:"section>entry>organizer>component>act>text"`
	EntryRelationshipOrganizer []organizer `xml:"section>entry>act>entryRelationship>organizer"`
}
