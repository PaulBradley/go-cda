package cda

type healthCareFacility struct {
	Code    []code `xml:"code"`
	Name    string `xml:"location>name"`
	Address []addr `xml:"location>addr"`
}

type encompassingEncounter struct {
	Code               []code               `xml:"code"`
	AssignedPerson     []person             `xml:"encounterParticipant>assignedEntity>assignedPerson>name"`
	HealthCareFacility []healthCareFacility `xml:"location>healthCareFacility"`
}
