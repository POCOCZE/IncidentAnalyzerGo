package main

import "testing"

func TestMTTRSec(t *testing.T) {
	report := &IncidentReport{
		IncidentsCount: 4,
        UnresolvedIDs: []string{"Incident1", "Incident2"},
        MTTR: "",
	}

    private := &IncidentPrivate{
        Duration: map[string]IncidentDuration{
            "Incident1": {Seconds: 10, HMSFormat: "0m10s"},
            "Incident2": {Seconds: 10, HMSFormat: "0m10s"},
        },
    }

    report.CalcMTTRSec(private)
    got := report.MTTR
    // 20 / 2 = 10s as MTTR
    want := "10s"
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}