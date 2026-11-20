package main

import (
	"strings"
	"testing"
)

func TestReadTimelineSortsEventsAndPreservesTies(t *testing.T) {
	input := "2026-11-20T09:05:00Z,Health check failed\n" +
		"2026-11-20T14:30:00+05:30,\"Deployment completed, revision 42\"\n" +
		"2026-11-20T09:00:00Z,   Readiness check started   \n"

	events, err := readTimeline(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"Deployment completed, revision 42", "Readiness check started", "Health check failed"}
	if len(events) != len(want) {
		t.Fatalf("got %d events, want %d", len(events), len(want))
	}
	for i, message := range want {
		if events[i].message != message {
			t.Errorf("event %d = %q, want %q", i, events[i].message, message)
		}
	}
}

func TestReadTimelineReportsInvalidTimestamp(t *testing.T) {
	input := "2026-11-20T09:00:00Z,Deployment completed\n" +
		"invalid,Health check failed\n"

	events, err := readTimeline(strings.NewReader(input))
	if err == nil || !strings.Contains(err.Error(), "record 2 timestamp") {
		t.Fatalf("got error %v, want a record 2 timestamp error", err)
	}
	if len(events) != 0 {
		t.Fatal("malformed input returned a partial timeline")
	}
}
