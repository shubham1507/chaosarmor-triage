package main

import "fmt"

type IncidentSnapshot struct {
	Service          string
	FailedRequests   int
	RecentDeployment bool
}

func (snapshot IncidentSnapshot) print() {
	fmt.Printf("Service: %s\n", snapshot.Service)
	fmt.Printf("Failed requests: %d\n", snapshot.FailedRequests)
	fmt.Printf("Recent deployment: %t\n", snapshot.RecentDeployment)
}
