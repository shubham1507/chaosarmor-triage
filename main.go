package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const defaultService = "unknown"

func main() {
	serviceName := defaultService
	var failedRequestsText string
	var failedRequests int
	var deployed bool
	var showTypes bool

	flag.StringVar(&serviceName, "service", defaultService, "Affected service name")
	flag.StringVar(&failedRequestsText, "failed-requests", "", "Number of failed requests")
	flag.BoolVar(&deployed, "deployed", false, "A deployment just happened")
	flag.BoolVar(&showTypes, "show-types", false, "Print snapshot field types")
	flag.Parse()

	if failedRequestsText != "" {
		if parsedCount, err := strconv.Atoi(failedRequestsText); err != nil || parsedCount < 0 {
			fmt.Fprintln(os.Stderr, "failed-requests must be a non-negative integer")
			os.Exit(2)
		} else {
			failedRequests = parsedCount
		}
	}

	fmt.Printf("Service: %s\n", serviceName)
	fmt.Printf("Failed requests: %d\n", failedRequests)
	fmt.Printf("Recent deployment: %t\n", deployed)

	if showTypes {
		fmt.Printf("Field types: %T, %T, %T\n", serviceName, failedRequests, deployed)
	}
}
