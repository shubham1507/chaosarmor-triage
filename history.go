package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseFailureWindows(history string, currentFailures, totalRequests int) ([]int, error) {
	entries := strings.Split(history, ",")
	windows := make([]int, 0, len(entries)+1)

	for index, entry := range entries {
		count, err := strconv.Atoi(strings.TrimSpace(entry))
		if err != nil || count < 0 || count > totalRequests {
			return nil, fmt.Errorf("history window %d must contain an integer from 0 to %d", index+1, totalRequests)
		}
		windows = append(windows, count)
	}

	return append(windows, currentFailures), nil
}

func forEachTransition(windows []int, visit func(previous, current int)) {
	for index := 1; index < len(windows); index++ {
		visit(windows[index-1], windows[index])
	}
}

func summarizeFailureTrends(windows []int) map[string]int {
	counts := make(map[string]int)

	forEachTransition(windows, func(previous, current int) {
		trend := "unchanged"
		if current < previous {
			trend = "improved"
		} else if current > previous {
			trend = "worsened"
		}
		counts[trend]++
	})

	return counts
}

func printFailureTrends(windows []int) {
	trendOrder := [...]string{"improved", "unchanged", "worsened"}
	counts := summarizeFailureTrends(windows)

	fmt.Printf("Failure counts (oldest to newest): %v\n", windows)
	for _, trend := range trendOrder {
		fmt.Printf("  %s: %d\n", trend, counts[trend])
	}
}
