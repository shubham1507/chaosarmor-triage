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

func printFailureTrends(windows []int) {
	trendOrder := [...]string{"improved", "unchanged", "worsened"}
	counts := make(map[string]int)

	for index, current := range windows[1:] {
		previous := windows[index]
		trend := "unchanged"
		if current < previous {
			trend = "improved"
		} else if current > previous {
			trend = "worsened"
		}
		counts[trend]++
	}

	fmt.Printf("Failure counts (oldest to newest): %v\n", windows)
	for _, trend := range trendOrder {
		fmt.Printf("  %s: %d\n", trend, counts[trend])
	}
}
