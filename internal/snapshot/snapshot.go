// Package snapshot models and writes incident snapshots.
package snapshot

import (
	"fmt"
	"io"
)

// IncidentSnapshot records request failures and deployment status for a service.
type IncidentSnapshot struct {
	Service          string
	FailedRequests   int
	RecentDeployment bool
}

var (
	_ io.WriterTo = IncidentSnapshot{}
	_ io.WriterTo = (*IncidentSnapshot)(nil)
)

// WriteTo writes the snapshot as text and returns the number of bytes written.
func (snapshot IncidentSnapshot) WriteTo(w io.Writer) (int64, error) {
	written, err := fmt.Fprintf(w, "Service: %s\nFailed requests: %d\nRecent deployment: %t\n",
		snapshot.Service, snapshot.FailedRequests, snapshot.RecentDeployment)
	return int64(written), err
}
