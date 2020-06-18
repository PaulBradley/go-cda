package cda

type organizer struct {
	ID            []id            `xml:"id"`
	StatusCode    []code          `xml:"statusCode"`
	EffectiveTime []effectiveTime `xml:"effectiveTime"`
}
