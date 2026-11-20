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
	totalRequestsText := flag.String("total-requests", "", "Total requests in the same observation window")
	flag.Parse()

	var totalRequests int
	if *totalRequestsText != "" {
		parsedTotal, err := strconv.Atoi(*totalRequestsText)
		if err != nil || parsedTotal <= 0 {
			fmt.Fprintln(os.Stderr, "total-requests must be a positive integer")
			os.Exit(2)
		}
		totalRequests = parsedTotal
	}

	if failedRequestsText != "" {
		if parsedCount, err := strconv.Atoi(failedRequestsText); err != nil || parsedCount < 0 {
			fmt.Fprintln(os.Stderr, "failed-requests must be a non-negative integer")
			os.Exit(2)
		} else {
			failedRequests = parsedCount
		}
	}

	if totalRequests > 0 && failedRequests > totalRequests {
		fmt.Fprintln(os.Stderr, "failed-requests cannot exceed total-requests")
		os.Exit(2)
	}

	fmt.Printf("Service: %s\n", serviceName)
	fmt.Printf("Failed requests: %d\n", failedRequests)
	fmt.Printf("Recent deployment: %t\n", deployed)

	if totalRequests > 0 {
		printAssessment(failedRequests, totalRequests, deployed)
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
