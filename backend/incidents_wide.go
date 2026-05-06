package main

func IncidentsWide(incidents []Incident) ([]Incident, error) {
	var incidentsWide []Incident
	durations := make(map[string]IncidentDuration)
	for _, incident := range incidents {
		calcIncidentDuration(incident, durations)
		message, isResolved := IncidentMessage(incident, durations)

		// Set gathered variables to each incident
		incident.Message = message
		incident.IsResolved = isResolved
		// Add this updated (wide) incident to slice
		incidentsWide = append(incidentsWide, incident)
	}

	return incidentsWide, nil
}