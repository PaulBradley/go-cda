package cda

type organizer struct {
	ID            []id            `xml:"id"`
	StausCode     []code          `xml:"statusCode"`
	EffectiveTime []effectiveTime `xml:"effectiveTime"`
}
