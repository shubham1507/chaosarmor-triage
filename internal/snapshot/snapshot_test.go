package snapshot

import (
	"bytes"
	"io"
	"testing"
)

func TestIncidentSnapshotWriteTo(t *testing.T) {
	snapshot := IncidentSnapshot{
		Service:          "checkout",
		FailedRequests:   12,
		RecentDeployment: true,
	}

	var report io.WriterTo = snapshot
	var output bytes.Buffer
	written, err := report.WriteTo(&output)
	if err != nil {
		t.Fatalf("WriteTo returned an error: %v", err)
	}

	want := "Service: checkout\nFailed requests: 12\nRecent deployment: true\n"
	if got := output.String(); got != want {
		t.Errorf("snapshot output = %q, want %q", got, want)
	}
	if written != int64(len(want)) {
		t.Errorf("bytes written = %d, want %d", written, len(want))
	}
}
