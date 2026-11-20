package main

import (
	"github.com/shubham1507/chaosarmor-triage/internal/snapshot"

	"encoding/json"
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
	totalRequestsText := flag.String("total-requests", "", "Total requests in the same observation window")
	historyText := flag.String("history", "", "Earlier failed-request counts, oldest first, comma-separated")
	snapshotJSON := flag.Bool("snapshot-json", false, "Print the basic incident snapshot as JSON")
	flag.Parse()

	totalRequests, err := parseOptionalTotal(*totalRequestsText)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	if failedRequestsText != "" {
		if parsedCount, err := strconv.Atoi(failedRequestsText); err != nil || parsedCount < 0 {
			fmt.Fprintln(os.Stderr, "failed-requests must be a non-negative integer")
			os.Exit(2)
		} else {
			failedRequests = parsedCount
		}
	}

	if totalRequests != nil && failedRequests > *totalRequests {
		fmt.Fprintln(os.Stderr, "failed-requests cannot exceed total-requests")
		os.Exit(2)
	}

	var windows []int
	if *historyText != "" {
		if totalRequests == nil {
			fmt.Fprintln(os.Stderr, "history requires total-requests")
			os.Exit(2)
		}

		var err error
		windows, err = parseFailureWindows(*historyText, failedRequests, *totalRequests)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}

	incident := snapshot.IncidentSnapshot{
		Service:          serviceName,
		FailedRequests:   failedRequests,
		RecentDeployment: deployed,
	}
	if *snapshotJSON {
		if err := json.NewEncoder(os.Stdout).Encode(incident); err != nil {
			fmt.Fprintln(os.Stderr, "could not encode incident snapshot:", err)
			os.Exit(1)
		}
		return
	}
	if _, err := incident.WriteTo(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "could not write incident snapshot:", err)
		os.Exit(1)
	}

	if totalRequests != nil {
		printAssessment(failedRequests, *totalRequests, deployed)
	}

	if len(windows) > 0 {
		printFailureTrends(windows)
	}

	if showTypes {
		fmt.Printf("Field types: %T, %T, %T\n", serviceName, failedRequests, deployed)
	}
}

func printAssessment(failedRequests, totalRequests int, deployed bool) {
	const (
		warningRate  = 5.0
		criticalRate = 10.0
	)

	failureRate := float64(failedRequests) / float64(totalRequests) * 100

	var severity string
	if failureRate >= criticalRate {
		severity = "critical"
	} else if failureRate >= warningRate {
		severity = "warning"
	} else {
		severity = "normal"
	}

	const (
		deploymentSignal = 1 << iota
		elevatedRateSignal
	)
	var signals uint
	if deployed {
		signals |= deploymentSignal
	}
	if failureRate >= warningRate {
		signals |= elevatedRateSignal
	}

	fmt.Printf("Failure rate: %.2f%%\n", failureRate)
	fmt.Printf("Severity: %s\n", severity)

	fmt.Print("Signals: ")
	printedSignals := 0
	for index, label := range []string{"recent deployment", "elevated failure rate"} {
		active := signals&(1<<index) != 0
		if active && printedSignals > 0 {
			fmt.Print(", ")
		}
		if active {
			fmt.Print(label)
			printedSignals++
		}
	}
	if printedSignals == 0 {
		fmt.Print("none")
	}
	fmt.Println()

	switch severity {
	case "critical":
		fmt.Println("Next check: inspect live errors and verify the scope of impact")
	case "warning":
		fmt.Println("Next check: compare errors before and after the latest change")
	default:
		fmt.Println("Next check: continue monitoring the service")
	}
}

func parseOptionalTotal(text string) (*int, error) {
	if text == "" {
		return nil, nil
	}

	total, err := strconv.Atoi(text)
	if err != nil || total <= 0 {
		return nil, fmt.Errorf("total-requests must be a positive integer")
	}

	return &total, nil
}
