package main

import (
	"fmt"
	"io"
)

type IncidentSnapshot struct {
	Service          string
	FailedRequests   int
	RecentDeployment bool
}

var (
	_ io.WriterTo = IncidentSnapshot{}
	_ io.WriterTo = (*IncidentSnapshot)(nil)
)

func (snapshot IncidentSnapshot) WriteTo(w io.Writer) (int64, error) {
	written, err := fmt.Fprintf(w, "Service: %s\nFailed requests: %d\nRecent deployment: %t\n",
		snapshot.Service, snapshot.FailedRequests, snapshot.RecentDeployment)
	return int64(written), err
}
